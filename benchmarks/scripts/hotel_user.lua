local function get_user()
  local id = math.random(0, 10000)
  local user_name = "Cornell_" .. tostring(id)
  local pass_word = ""
  for i = 0, 9, 1 do 
    pass_word = pass_word .. tostring(id)
  end
  return user_name, pass_word
end

local function user_login()
  local user_name, password = get_user()
  local method = "POST"
  local path = "http://localhost:8080/r/hotel/user"
  local headers = {["Content-Type"] = "application/json"}
  local body = string.format('{"username":"%s","password":"%s"}', user_name, password)
  return wrk.format(method, path, headers, body)
end

request = function()
  return user_login()
end

function init()
  rand_seed = os.time()
  math.randomseed(rand_seed)
end
