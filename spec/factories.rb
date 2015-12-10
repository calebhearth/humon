Event = Class.new(OpenStruct) do
  def reload
    self
  end
end
User = Class.new(OpenStruct)

FactoryGirl.define do
  sequence :lat do |n|
    "#{n}.0".to_f
  end

  sequence :lon do |n|
    "#{n}.0".to_f
  end

  sequence :name do |n|
    "name #{n}"
  end

  sequence :started_at do
    DateTime.now
  end

  sequence :token do
    SecureRandom.hex(3)
  end

  factory :attendance do
    event
    user
  end

  factory :event do
    sequence(:id)
    lat
    lon
    name
    started_at
    owner factory: :user
  end

  factory :user do
    device_token { generate(:token) }
  end
end
