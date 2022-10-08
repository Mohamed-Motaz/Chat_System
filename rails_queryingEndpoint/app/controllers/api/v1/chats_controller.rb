module Api
    module V1
        class ChatsController < ApplicationController
            # GET /applications/:token/chats
            def index
                @token = (params[:token] || " ").strip
                if !@token || @token.empty?
                    puts @token
                    return render json: {error: "Invalid token entered"} , status: :bad_request
                end

                @res = Chat.find_by_sql("
                    SELECT number, messages_count 
                    FROM instabug.chats
                    WHERE application_token = '#{@token}'
                    ")

                if @res.length > 0
                    render json: @res.as_json(only: [:number, :messages_count]), status: :ok
                else
                    render json: {error: "Invalid token entered"} , status: :bad_request
                end
            end

            # GET /applications/:token/chats/:chat_num
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


                @res = Chat.find_by_sql("
                    SELECT number, messages_count 
                    FROM instabug.chats
                    WHERE application_token = '#{@token}' and number = #{@chat_num}
                    LIMIT 1")

                if @res.length > 0
                    render json: @res[0].as_json(only: [:number, :messages_count]), status: :ok
                else
                    render json: {error: "Invalid token entered"} , status: :bad_request
                end
            end

        end
    end
end