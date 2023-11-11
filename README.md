# Project: Data Engineering

This is the monorepo containing all necessary components to run the project for the data engineering course.

## Service Descriptions

- **Instagram API Mock Service:** This services continuously pushes data into the Kafka topics, to simulate the process of scraping data from the Instagram API. Additionally the service creates mock pictures by calling thispersondoesnotexist.com and hosts them under the path /images.
- **Storage Service:** This service subscribes to the Kafka Topic: `instagram-profiles` and upon new messages creates new entries in the postgresql data base. Any new media entries it finds (e.g. posts) triggers a download of the associated image, that is then uploaded to the minio instance.
- **Data Access Service:** This REST-service can be called at `/api/all` to get a basic aggregation of the insight data of all currently tracked profiles.
- **Visualization Service:** This service hosts a nodejs-based frontend that visualizes some basic statistics (number of likes, reach, saved, comments) over the lifetime of all media.

## Configuration

All configuration should be done in the accompanying docker compose file using the environment variables.  
All variables need to be set in the docker-compose file or the project won't start.  


| Variable                 	| Description                                                                       	| Default             	|
|--------------------------	|-----------------------------------------------------------------------------------	|---------------------	|
| S3_EXTERNAL_ENDPOINT     	| The external endpoint under which the S3-based storage can be accessed            	| localhost:10002     	|
| S3_ENDPOINT              	| The docker-network internal endpoint to call the S3-storage                       	| minio:10002         	|
| S3_ACCESS_KEY_ID         	| Access key for accessing the minio storage/ui as a root user                      	| admin               	|
| S3_SECRET_ACCESS_KEY     	| The password for accessing the minio storage/ui as a root user                    	| admin123            	|
| POSTGRES_USER            	| Default postgres database user                                                    	| postgres            	|
| POSTGRES_PASSWORD        	| Default password for the POSTGRES_USER                                            	| test                	|
| POSTGRES_DB              	| Default database to create/use in postgres                                        	| postgres            	|
| POSTGRES_HOST            	| Host under which postgres is reachable                                            	| postgres            	|
| POSTGRES_PORT            	| Port used to communicate with postgres                                            	| 5432                	|
| CLICKHOUSE_ENDPOINT      	| Endpoint used by the internal clients to communicate with the clickhouse instance 	| clickhouse:9000     	|
| DATA_ACCESS_SERVICE_HOST 	| Host called from within the visualization service to receive data                 	| data-access-service 	|
| DATA_ACCESS_SERVICE_PORT 	| Port called from within the visualization service to receive data                 	| 3001                	|
| KAFKA_BOOTSTRAP_SERVERS  	| Location of kafka servers                                                         	| kafka:9092          	|


Additionally the instagram API mock service offers a few additional variables to configure the data creation.

| Variable                 	| Description                                                       	| Default 	|
|--------------------------	|-------------------------------------------------------------------	|---------	|
| NUM_PROFILES             	| Number of profiles to simulate                                    	| 10      	|
| NUM_PICTURES_PER_PROFILE 	| Number of pictures per profile to simulate                        	| 10      	|
| INSIGHT_UPDATE_FREQ_MS   	| Time in milliseconds between updates for insights of a media post 	| 30      	|
| PROFILE_UPDATE_FREQ_MS   	| Time in milliseconds between updates for a singular profile       	| 50      	|

## Running the Project

There are currently two ways to start up the project:

1. Use the included reset.sh-script. This script deletes any old leftovers of previous runs, rebuilds all containers and starts the environment.
```bash
./reset.sh
```

2. Use docker compose directly to start the application

```bash
docker compose build 
docker compose up -d 
```

The project can be configured by adjusting the variables in the docker compose file.  
Depending on the variables, the required start up time can vary in length. The Instagram API Mock Service in particular requires a longer start time, as there is a 500ms delay between fetching the individual images from thispersondoesnotexist.com to guarantee that unique images are fetched.  
Therefore, this service usually needs (NUM_PROFILES*NUM_PICTURES_PER_PROFILE / 2) seconds to start.  
After starting the project the data can be visualized by opening the url: http://localhost:5473 in a browser.

