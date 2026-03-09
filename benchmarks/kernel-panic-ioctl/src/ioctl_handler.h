#ifndef IOCTL_HANDLER_H
#define IOCTL_HANDLER_H

#include "device.h"

/* ioctl command codes */
#define IOCTL_SET_CONFIG  1
#define IOCTL_GET_CONFIG  2
#define IOCTL_RESET       3
#define IOCTL_GET_STATUS  4

/* Request structure — in a real kernel driver, the value/key pointers
 * would come from userspace and need copy_from_user/copy_to_user.
 * Here we simulate the same pattern: the handler receives pointers
 * that may be untrusted (NULL, invalid, etc.) */
struct ioctl_request {
    int cmd;
    const char *key;
    char *value;
    size_t value_len;
};

/* Handle an ioctl request on the given device.
 * Returns 0 on success, -1 on error. */
int handle_ioctl(struct device *dev, struct ioctl_request *req);

#endif /* IOCTL_HANDLER_H */
