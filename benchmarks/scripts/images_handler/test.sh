wrk --latency -R 50 http://pc70.cloudlab.umass.edu:8080/r/hotel/user > ./hotel_user/optimized/50
wrk --latency -R 100 http://pc70.cloudlab.umass.edu:8080/r/hotel/user > ./hotel_user/optimized/100
wrk --latency -R 100 http://pc70.cloudlab.umass.edu:8080/r/hotel/user > ./hotel_user/optimized/100
wrk --latency -R 200 http://pc70.cloudlab.umass.edu:8080/r/hotel/user > ./hotel_user/optimized/200
wrk --latency -R 300 http://pc70.cloudlab.umass.edu:8080/r/hotel/user > ./hotel_user/optimized/300
wrk --latency -R 400 http://pc70.cloudlab.umass.edu:8080/r/hotel/user > ./hotel_user/optimized/400
wrk --latency -R 500 http://pc70.cloudlab.umass.edu:8080/r/hotel/user > ./hotel_user/optimized/500


cat ./hotel_user/100 ./hotel_user/200 ./hotel_user/300 ./hotel_user/400 | ./venv/bin/wrk2img ./hotel_user/output.png 


# original
wrk --latency -R 50 http://pc70.cloudlab.umass.edu:8080/r/app/hotels > ./hotel_user/original/50


# Optimized, qps
cat ./hotel_user/optimized/50 ./hotel_user/optimized/100 ./hotel_user/optimized/200 ./hotel_user/optimized/300 ./hotel_user/optimized/400 | ./venv/bin/wrk2img output.png 


cat ./hotel_user/original/50 ./hotel_user/optimized/50 | ./venv/bin/wrk2img original_optimized.png 