require 'httparty'

KEY = "8fc630dc5d5fdbb723d097b833d8fb93"
TOKEN = "0fc1f3a058d4c08dde0da8ed2d51618caf37f13a4bf671b0b44e28256a32033e"

class Trello
  include HTTParty
  base_uri "https://api.trello.com/1"

  def with_auth(options)
    options.merge(key: KEY, token: TOKEN)
  end

  def new_board(query={})
    JSON.parse(self.class.post("/boards", query: with_auth(query)).body)
  end

  def new_list(query={})
    JSON.parse(self.class.post("/lists", query: with_auth(query)).body)
  end

  def new_card(query={})
    JSON.parse(self.class.post("/cards", query: with_auth(query)).body)
  end
end

trello = Trello.new

json = JSON.parse(File.read("entities.json")) #this is based on the JSON generated from the corresponding "trello.lang"

json["Board"].each do |b|
  board = trello.new_board(name: b["name"])
  b["lists"].each do |l|
    list = trello.new_list(name: l["name"], idBoard: board["id"])
    l["cards"].each do |c|
      trello.new_card(name: c["name"], idList: list["id"], idBoard: board["id"])
    end
  end
end
