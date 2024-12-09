local function get_user()
  local id = math.random(0, 10000)
  local user_name = "Cornell_" .. tostring(id)
  local pass_word = ""
  for i = 0, 9, 1 do 
    pass_word = pass_word .. tostring(id)
  end
  return user_name, pass_word
end

local function search_hotel() 
  local in_date = math.random(9, 23)
  local out_date = math.random(in_date + 1, 24)

  local in_date_str = tostring(in_date)
  if in_date <= 9 then
    in_date_str = "2015-04-0" .. in_date_str 
  else
    in_date_str = "2015-04-" .. in_date_str
  end

  local out_date_str = tostring(out_date)
  if out_date <= 9 then
    out_date_str = "2015-04-0" .. out_date_str 
  else
    out_date_str = "2015-04-" .. out_date_str
  end

  local lat = 38.0235 + (math.random(0, 481) - 240.5)/1000.0
  local lon = -122.095 + (math.random(0, 325) - 157.0)/1000.0

  local method = "POST"
  local path = "http://localhost:8080/r/hotel/hotels"
  local headers = {["Content-Type"] = "application/json"}
  local body = string.format('{"inDate":"%s","outDate":"%s","lat":%f,"lon":%f}', in_date_str, out_date_str, lat, lon)
  return wrk.format(method, path, headers, body)
end

local function recommend()
  local coin = math.random()
  local req_param = ""
  if coin < 0.33 then
    req_param = "dis"
  elseif coin < 0.66 then
    req_param = "rate"
  else
    req_param = "price"
  end

  local lat = 38.0235 + (math.random(0, 481) - 240.5)/1000.0
  local lon = -122.095 + (math.random(0, 325) - 157.0)/1000.0

  local method = "POST"
  local path = "http://localhost:8080/r/hotel/recommendations"
  local headers = {["Content-Type"] = "application/json"}
  local body = string.format('{"require":"%s","lat":%f,"lon":%f}', req_param, lat, lon)
  return wrk.format(method, path, headers, body)
end

local function reserve()
  local in_date = math.random(9, 23)
  local out_date = in_date + math.random(1, 5)

  local in_date_str = tostring(in_date)
  if in_date <= 9 then
    in_date_str = "2015-04-0" .. in_date_str 
  else
    in_date_str = "2015-04-" .. in_date_str
  end

  local out_date_str = tostring(out_date)
  if out_date <= 9 then
    out_date_str = "2015-04-0" .. out_date_str 
  else
    out_date_str = "2015-04-" .. out_date_str
  end

  local hotel_id = tostring(math.random(1, 1000))
  local user_id, password = get_user()
  local cust_name = user_id

  local num_room = "1"

  local method = "POST"
  local path = "http://localhost:8080/r/hotel/reservation"
  local headers = {["Content-Type"] = "application/json"}
  local body = string.format('{"inDate":"%s","outDate":"%s","hotelId":"%s","customerName":"%s","username":"%s","password":"%s","number":"%s"}', in_date_str, out_date_str, hotel_id, cust_name, user_id, password, num_room)
  return wrk.format(method, path, headers, body)
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
  local search_ratio      = 0.6
  local recommend_ratio   = 0.39
  local user_ratio        = 0.005
  local reserve_ratio     = 0.005

  local coin = math.random()
  if coin < search_ratio then
    return search_hotel()
  elseif coin < search_ratio + recommend_ratio then
    return recommend()
  elseif coin < search_ratio + recommend_ratio + user_ratio then
    return user_login()
  else 
    return reserve()
  end
end

function init()
  rand_seed = os.time()
  math.randomseed(rand_seed)
end
