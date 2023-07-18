-- Creates Table in Clickhouse that later on contains the data from kafka
CREATE TABLE ig_insights 
(
	id UUID,
	comments UInt64,
	engagement UInt64,
	impressions UInt64,
	likes UInt64,
	reach UInt64,
	saved UInt64
) ENGINE = MergeTree ORDER BY (id);

-- Creates connector to kafka that consumes the messages and temporarily stores them in clickhouse
CREATE TABLE ig_insights_queue
(
	id UUID,
	comments UInt64,
	engagement UInt64,
	impressions UInt64,
	likes UInt64,
	reach UInt64,
	saved UInt64
) ENGINE = Kafka settings 
kafka_broker_list= 'kafka:9092', 
kafka_topic_list = 'instagram-insights', 
kafka_group_name = 'clickhouse',
kafka_format = 'JSONEachRow',
kafka_thread_per_consumer = 0, 
kafka_num_consumers = 1;

-- Connector between the queue and the final table
CREATE MATERIALIZED VIEW default.ig_insights_mv TO default.ig_insights AS SELECT * FROM `default`.ig_insights_queue;

-- Clear table
-- DELETE FROM ig_insights;

-- Stop Kafka ingest
-- DETACH TABLE ig_insights_queue;

-- Restart Kafka Ingest
-- ATTACH TABLE ig_insights_queue;
