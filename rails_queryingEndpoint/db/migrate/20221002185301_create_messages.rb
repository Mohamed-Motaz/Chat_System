class CreateMessages < ActiveRecord::Migration[5.2]
  def change
    create_table :messages do |t|
      t.references :chat, foreign_key: true, null: false
      t.integer :number, null: false
      t.string :body
      t.index [:chat_id, :number]
      t.timestamps
    end
    #add_index :messages, [:chat_id, :number], :name => "index_messages_on_chat_id_and_number"

  end
end
