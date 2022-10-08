module Api
    module V1
        class MessagesController < ApplicationController
            # GET /applications/:token/chats/:chat_num/messages
            def index
                @token = (params[:token] || " ").strip
                if !@token || @token.empty?
                    puts @token
                    return render json: {error: "Invalid token entered"} , status: :bad_request
                end
                @chat_num = (params[:chat_num] || " ").strip
                if !@chat_num || @chat_num.empty?
                    puts @chat_num
                    return render json: {error: "Invalid chat num entered"} , status: :bad_request
                end

                @res = Message.find_by_sql("
                    SELECT instabug.messages.number, instabug.messages.body
                    FROM instabug.chats
                    INNER JOIN instabug.messages ON instabug.messages.chat_id = instabug.chats.id
                    WHERE instabug.chats.application_token = '#{@token}' and instabug.chats.number = #{@chat_num};
                    ")

                if @res.length > 0
                    render json: @res.as_json(only: [:number, :body]), status: :ok
                else
                    render json: {error: "Invalid token entered"} , status: :bad_request
                end
            end

            # GET /applications/:token/chats/:chat_num/messages/:message_num
            def show
                @token = (params[:token] || " ").strip
                if !@token || @token.empty?
                    puts @token
                    return render json: {error: "Invalid token entered"} , status: :bad_request
                end
                @chat_num = (params[:chat_num] || " ").strip
                if !@chat_num || @chat_num.empty?
                    puts @chat_num
                    return render json: {error: "Invalid chat num entered"} , status: :bad_request
                end
                @message_num = (params[:message_num] || " ").strip
                if !@message_num || @message_num.empty?
                    puts @message_num
                    return render json: {error: "Invalid message num entered"} , status: :bad_request
                end


                @res = Chat.find_by_sql("
                    SELECT instabug.messages.number, instabug.messages.body
                    FROM instabug.chats
                    INNER JOIN instabug.messages ON instabug.messages.chat_id = instabug.chats.id and instabug.messages.number = #{@message_num}
                    WHERE instabug.chats.application_token = '#{@token}' and instabug.chats.number = #{@chat_num}
                    LIMIT 1")

                if @res.length > 0
                    render json: @res[0].as_json(only: [:number, :body]), status: :ok
                else
                    render json: {error: "Invalid token entered"} , status: :bad_request
                end
            end

        end
    end
end