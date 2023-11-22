require 'json'
require 'digest/sha2'
require 'securerandom'
require 'json'
require "sqlite3"

key_cardinality = 10000
key_salt = 'toy-redis'
value_salts = ('a'...'z')

value_max_len = 5_000
value_min_len = 10

schema = {
  data: []
}

digest_key = Digest::SHA256.new

(1...key_cardinality).each do |kmc|
  digest_key.update(key_salt+kmc.to_s)
  key = digest_key.hexdigest  


  v_len = rand(value_min_len..value_max_len)
  value = SecureRandom.alphanumeric(v_len)

  schema[:data].push({
    key: key,
    value: value
  })
end

File.open("./loadtesting/testdata.json", "w") do |f|
  f.write(JSON.generate(schema))
end
