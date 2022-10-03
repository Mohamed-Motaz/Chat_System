Pic of system

Assumptions

Issues:
Memory may run out of redis before pushing the changes to the db
Redis is a single point of failure
Workers may die after pulling a job from the queue, so the job may be lost
Uuid token may actually collide with another (0.00001% chance) so it needs to be   handled.
Before accepting a message, I first need to check if the chat is present in redis or the db
A very complex design decision is maintaining the counter. The corner case that drove me nuts is as follows: We have a few go servers running, with the chats_number = 100. The value is currently stored in redis. Redis dies, and then the servers die as well. The only place that knows what the max number is is the db. The go servers are then woken up along with Redis. At this point, Redis doesn’t have a key called chats_number. So the solution is, for every go server, when it encounters a new chat that it hasn’t served before, iT first queries the db for the maximum chats_number. It then attempts to set the key chats_number in Redis to the correct value. After doing so, it calls INCR to get the new number. Here is the race condition: Two servers do the same sequence of events. A server “A” calls MySQL, while the other “B” has called MySQL, set the ctr, and incremented it appropriately. Why did both call MySQL? Because this is the first time both have seen this chat for this app. So they aren’t sure if it is in Redis or not. The issue is that now server “A” attempts to set the ctr to the value it obtained from MySQL. This is wrong, because the ctr has already been set and incremented. This would result in duplicate numbers. So the solution is to use SETNX (if key doesn’t exist, then calling INCR to get the appropriate counter). In essence, whenever a server receives a request to create a chat for the first time for a specific app, or a message for a specific chat for the first time, it has to call MySQL, SETNX the key, and then INCR it to get the correct result.

In the chats table, I chose the foreign key to be on the application_token, rather than the application_id. This was a hard choice. Pros: Only 1 query to get the applications chat. Cons: Indexing a varchar, which means that the index may grow in size quite a bit in the future. Comparing varchars is indeed slower than ints, but if I opt for using the application_id, the first query would still have to compare tokens to use the index. So it makes more sense to decrease the queries sent to the db which may already be under alot of load due to the messages.



Schema

Applications				//persisted immediately
Id: int
Token: uuid  //unique index
Name: string
Chats_count: int



Chats					//may lag for up to an hour
Id: int
applicationId: int  //foreign key
Number: int
Messages_count: int

Messages				//may lag for up to an hour
Id: int
chatId: int
Number: int
Body: text


