FROM mongo:3.4

COPY data /tmp/dump

CMD mongod --fork --logpath /var/log/mongodb.log; \
	mongorestore -d ${MONGO_INITDB_DATABASE} /tmp/dump/; \
	mongod --shutdown; \
	docker-entrypoint.sh mongod
