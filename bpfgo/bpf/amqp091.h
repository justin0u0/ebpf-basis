#ifndef __AMQP091_H__
#define __AMQP091_H__

#define AMQP_PORT 5672

// AMQP frame detail
struct amqp_frame_header {
	uint8_t type; // method, header, body, heartbeat
	uint16_t channel;
	uint32_t size;
} __attribute__((packed));

#define AMQP_FRAME_HEADER_SIZE 7
#define AMQP_FRAME_HEADER_TYPE_METHOD 1
#define AMQP_FRAME_HEADER_TYPE_HEADER 2
#define AMQP_FRAME_HEADER_TYPE_BODY 3
#define AMQP_FRAME_HEADER_TYPE_HEARTBEAT 8

// AMQP method frame

struct amqp_method_frame {
	uint16_t class_id;
	uint16_t method_id;
	// arguments
} __attribute__((packed));

#define AMQP_FRAME_END 0xCE

#endif
