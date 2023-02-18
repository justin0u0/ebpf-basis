#ifndef __AMQP091_H__
#define __AMQP091_H__

#define AMQP_PORT 5672

// AMQP_MAX_FRAME_SIZE and AMQP_MAX_NUM_FRAMES are used to limit the total bytes
// of AMQP frames that we process. Due to the limitation of the BPF verifier.
#define AMQP_MAX_FRAME_SIZE 4096
#define AMQP_MAX_NUM_FRAMES 8

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

// AMQP header frame

struct amqp_header_frame {
	uint16_t class_id;
	uint16_t weight;
	uint64_t body_size;
	// properties
} __attribute__((packed));

const uint8_t AMQP_FRAME_END = 0xCE;

#endif
