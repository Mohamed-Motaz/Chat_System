module Api
    module V1
        class ApplicationsController < ApplicationController

            # GET /applications
            def index
                @applications = Application.all
                render json: @applications.as_json(only: [:name, :token, :chats_count]), status: :ok
            end

            # GET /applications/:token
            def show
                @token = (params[:token] || " ").strip
                if !@token || @token.empty?
                    puts @token
                    render json: {error: "Invalid token entered"} , status: :bad_request
                else
                    @res = Application.find_by_sql("
                        SELECT name, token, chats_count 
                        FROM instabug.applications
                        WHERE token = '#{@token}'")
                    render json: @res.as_json(only: [:name, :token, :chats_count]), status: :ok
                end
            end

            # POST /applications
            def create
                @name = (params[:name] || " ").strip
                if !@name || @name.empty?
                    puts @name
                    render json: {error: "Invalid name entered"} , status: :bad_request
                else
                    @application = Application.new({:name => @name, :chats_count => 0, :token => SecureRandom.uuid})
                    if @application.save
                        render json: @application.as_json(only: [:name, :token, :chats_count]), status: :created
                    else
                        render json: @application.errors.as_json, status: :unprocessable_entity
                    end
                end
            end


            # PUT /applications/:token
            def update
                @name = (params[:name] || " ").strip
                @token = (params[:token] || " ").strip

                if !@name || @name.empty?
                    puts @name
                    render json: {error: "Invalid name entered"} , status: :bad_request
                elsif !@token || @token.empty?
                    puts @token
                    render json: {error: "Invalid token entered"} , status: :bad_request
                else
                    @res = Application.connection.exec_update(
                        Application.sanitize_sql([
                        "UPDATE instabug.applications 
					    SET name = :name
					    WHERE token = :token ",
                        {
                            name: @name,
                            token: @token
                        }
                        ]))
                    if @res == 1
                        render json: {success: "ok"}.as_json, status: :ok
                    else
                        render json: {error: "no such token"}.as_json, status: :bad_request
                    end
                end
            end

    
        end
    end
end