require "active_record"
ActiveRecord::Base.establish_connection(
  "postgres://localhost/humon_development"
)
class User < ActiveRecord::Base; end
class Event < ActiveRecord::Base
  belongs_to :user
end
class Attendance < ActiveRecord::Base
  belongs_to :user
  belongs_to :event
end

FactoryGirl.define do
  sequence :lat do |n|
    "#{n}.1".to_f
  end

  sequence :lon do |n|
    "#{n}.2".to_f
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
    address "address"
    lat
    lon
    name
    started_at
    ended_at { generate(:started_at) }
    user
  end

  factory :user do
    device_token { generate(:token) }
  end
end
