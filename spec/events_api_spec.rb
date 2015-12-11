require "spec_helper"
require "net/http"
require "json"
require "date"
require "factories"
require "pry"

describe "Events API" do
  include FactoryGirl::Syntax::Methods

  describe "GET /v1/events/:id" do
    it "404s for a bad ID" do
      response = get "/v1/events/9001"

      expect(response.code).to eq "404"
    end

    it "returns an event by :id" do
      event = create(:event)
      response = get "/v1/events/#{event.id}"
      expect(response.code).to eq "200"

      expect(JSON.parse(response.body)).to eq(
        "address" => event.address,
        "ended_at" => event.ended_at.utc.strftime('%Y-%m-%dT%H:%M:%S.%6NZ'),
        "id" => event.id,
        "lat" => event.lat,
        "lon" => event.lon,
        "name" => event.name,
        "owner" => { "id" => event.user_id },
        "started_at" => event.started_at.utc.strftime('%Y-%m-%dT%H:%M:%S.%6NZ'),
      )
    end

    def get(path)
      uri = URI("http://localhost:4321" + path)
      Net::HTTP.get_response(uri)
    end
  end

  describe 'POST /v1/events' do

    it 'saves the address, lat, lon, name, owner, and started_at date' do
      date = DateTime.now.utc
      datestring = date.strftime('%Y-%m-%dT%H:%M:%S.%6NZ')
      device_token = '123abcd456xyz'
      owner = create(:user, device_token: device_token)

      response = post '/v1/events', {
        address: '123 Example St.',
        ended_at: datestring,
        lat: 1.0,
        lon: 1.0,
        name: 'Fun Place!!',
        started_at: datestring,
        owner: { id: owner.id },
      }.to_json,
      set_headers(device_token)

      expect(response.code).to eq "200"
      response_json = JSON.parse(response.body)
      event = Event.last
      expect(response_json).to eq({ 'id' => event.id })
      expect(event.address).to eq '123 Example St.'
      expect(event.ended_at).to eq date
      expect(event.lat).to eq 1.0
      expect(event.lon).to eq 1.0
      expect(event.name).to eq 'Fun Place!!'
      expect(event.started_at).to eq date
      expect(event.user).to eq owner
    end

    it 'returns an error message when invalid' do
      device_token = '123abcd456xyz'

      response = post '/v1/events',
        {}.to_json,
        set_headers(device_token)

      response_json = JSON.parse(response.body)
      expect(response_json).to eq({
        'message' => 'Validation Failed',
        'errors' => [
          "Lat can't be blank",
          "Lon can't be blank",
          "Name can't be blank",
          "Started at can't be blank",
        ]
      })
      expect(response.code.to_i).to eq 422
    end

    def post(path, data, headers)
      uri = URI("http://localhost:4321" + path)
      Net::HTTP.start(uri.host, uri.port) do |http|
        http.post(uri, data, headers)
      end
    end
  end

  describe 'PATCH /v1/events/:id' do

    it 'updates the event attributes' do
      event = create(:event, name: 'Old name')
      new_name = 'New name'

      response = patch "/v1/events/#{event.id}", {
        address: event.address,
        ended_at: event.ended_at,
        lat: event.lat,
        lon: event.lon,
        name: new_name,
        started_at: event.started_at,
        owner: {
          id: event.user.id
        }
      }.to_json,
      set_headers(event.user.device_token)

      event = event.reload
      expect(event.name).to eq new_name
      expect(JSON.parse(response.body)).to eq({ 'id' => event.id })
    end

    it 'returns an error message when invalid' do
      event = create(:event)

      response = patch "/v1/events/#{event.id}", {
        address: event.address,
        ended_at: event.ended_at,
        lat: event.lat,
        lon: event.lon,
        name: nil,
        started_at: event.started_at,
        owner: {
          id: event.user.id
        }
      }.to_json,
      set_headers(event.user.device_token)

      event = event.reload
      expect(event.name).to_not be nil
      expect(JSON.parse(response.body)).to eq({
        'message' => 'Validation Failed',
        'errors' => [
          "Name can't be blank"
        ]
      })
      expect(response.code.to_i).to eq 422
    end

    def patch(path, data, headers)
      uri = URI("http://localhost:4321" + path)
      Net::HTTP.start(uri.host, uri.port) do |http|
        http.patch(uri, data, headers)
      end
    end
  end


  def set_headers(device_token)
    app_secret = 'secretkey'
    ENV['TB_APP_SECRET'] = app_secret

    {
      'tb-app-secret' => app_secret,
      'tb-device-token' => device_token,
      'Content-Type' => 'application/json'
    }
  end
end
