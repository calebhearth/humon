require "spec_helper"
require "net/http"
require "json"
require "date"

describe "GET /v1/events/:id" do
  it "404s for a bad ID" do
    response = get "/v1/events/9001"

    expect(response.code).to eq "404"
  end

  it "returns an event by :id" do
    response = get "/v1/events/1"
    expect(response.code).to eq "200"

    expect(JSON.parse(response.body)).to eq(
      "address" => "123 Main St",
      "ended_at" => "2001-01-01T00:00:00Z",
      "id" => 1,
      "lat" => "30.267153",
      "lon" => "-97.743061",
      "name" => "Austin",
      "owner" => { "id" => 1 },
      "started_at" => "2001-01-01T00:00:00Z",
    )
  end

  def get(path)
    uri = URI("http://localhost:4321" + path)
    Net::HTTP.start(uri.host, uri.port) do |http|
      request = Net::HTTP::Get.new uri
      return http.request request
    end
  end
end
