wrk --latency -R 100 http://pc70.cloudlab.umass.edu:8080/r/hotel/user > ./hotel_user/100
wrk --latency -R 200 http://pc70.cloudlab.umass.edu:8080/r/hotel/user > ./hotel_user/200
wrk --latency -R 300 http://pc70.cloudlab.umass.edu:8080/r/hotel/user > ./hotel_user/300
wrk --latency -R 400 http://pc70.cloudlab.umass.edu:8080/r/hotel/user > ./hotel_user/400
wrk --latency -R 500 http://pc70.cloudlab.umass.edu:8080/r/hotel/user > ./hotel_user/500


cat ./hotel_user/100 ./hotel_user/200 ./hotel_user/300 ./hotel_user/400 | ./venv/bin/wrk2img ./hotel_user/output.png 