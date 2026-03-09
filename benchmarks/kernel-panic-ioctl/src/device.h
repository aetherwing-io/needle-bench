#ifndef DEVICE_H
#define DEVICE_H

#include <stddef.h>

#define MAX_CONFIG_ENTRIES 64
#define MAX_KEY_LEN 64
#define MAX_VALUE_LEN 256
#define DEVICE_NAME_LEN 32

struct config_entry {
    char key[MAX_KEY_LEN];
    char value[MAX_VALUE_LEN];
    int in_use;
};

struct device {
    char name[DEVICE_NAME_LEN];
    struct config_entry config[MAX_CONFIG_ENTRIES];
    int config_count;
    int initialized;
    unsigned long flags;
};

struct device *device_create(const char *name);
void device_destroy(struct device *dev);
void device_print_info(struct device *dev);
int device_self_test(struct device *dev);
int device_fuzz_test(struct device *dev);

#endif /* DEVICE_H */
