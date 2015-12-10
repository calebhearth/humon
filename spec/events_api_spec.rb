require "spec_helper"
require "net/http"
require "json"

describe "GET /v1/events/:id" do
  it "returns an event by :id" do
    response = get "/v1/events/1"

    expect(response).to eq(
      "address" => "123 Main St",
      "ended_at" => "12:00 AM 1/1/2001",
      "id" => 1,
      "lat" => "30.267153",
      "lon" => "-97.743061",
      "name" => "Austin",
      "owner" => { "id" => 1 },
      "started_at" => "1/1/2001",
    )
  end

  def get(path)
    response = Net::HTTP.get(URI("http://localhost:4321" + path))
    JSON.parse(response.body)
  end
end
