docker run -d --name mongodb-hotelresv --hostname socialnetwork-mongodb -p 27017:27017 -v /mnt/inmem/db:/data/db mongo:4.2.8-bionic mongod --nojournal

docker run -d --name memcached-reserve --hostname memcached-reserve --restart always -m 4096m -p 11211:11211 memcached:1.5.22 memcached -m 4096