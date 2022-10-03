class Application < ApplicationRecord
    has_secure_token :token
    has_many :chats
    alias_attribute :id, :token



end
