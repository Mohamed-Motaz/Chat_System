class CreateChats < ActiveRecord::Migration[5.2]
  def change
    create_table :chats do |t|
      t.references :application, foreign_key: true, null: false
      t.integer :number, null: false
      t.integer :messages_count, null: false, default: 0
      t.index [:application_id, :number, :messages_count]
      t.timestamps
    end
    #add_index(:chats, [:application_id, :number, :messages_count]) #a covering index :)

  end
end
