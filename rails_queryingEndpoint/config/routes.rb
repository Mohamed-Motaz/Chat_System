Rails.application.routes.draw do
  # For details on the DSL available within this file, see http://guides.rubyonrails.org/routing.html
  namespace :api do
    namespace :v1 do
      get 'applications' => 'applications#index'
      get 'applications/:token' => 'applications#show'
      post 'applications' => 'applications#create'
      put 'applications/:token' => 'applications#update'

     
    end
  end
end