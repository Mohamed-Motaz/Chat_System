FROM ruby:2.7.6
RUN apt-get update -qq
RUN mkdir /app
WORKDIR /app
COPY ./Gemfile /app/Gemfile
COPY ./Gemfile.lock /app/Gemfile.lock
RUN gem install bundler:2.3.22
RUN bundle install
COPY . /app/.

EXPOSE 3000
CMD ["bin/rails", "server", "-b", "0.0.0.0"]

#docker build -t ruby_server_image .
#docker run -p 3000:3000 --name ruby_container ruby_server_image