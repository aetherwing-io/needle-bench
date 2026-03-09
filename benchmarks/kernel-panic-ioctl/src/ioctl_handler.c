#include <stdio.h>
#include <string.h>

#include "ioctl_handler.h"

/*
 * Find a config entry by key. Returns pointer to entry or NULL.
 */
static struct config_entry *find_config(struct device *dev, const char *key) {
    for (int i = 0; i < MAX_CONFIG_ENTRIES; i++) {
        if (dev->config[i].in_use && strcmp(dev->config[i].key, key) == 0) {
            return &dev->config[i];
        }
    }
    return NULL;
}

/*
 * Find a free config slot. Returns pointer to entry or NULL if full.
 */
static struct config_entry *find_free_slot(struct device *dev) {
    for (int i = 0; i < MAX_CONFIG_ENTRIES; i++) {
        if (!dev->config[i].in_use) {
            return &dev->config[i];
        }
    }
    return NULL;
}

/*
 * Handle IOCTL_SET_CONFIG: set a configuration key-value pair.
 */
static int handle_set_config(struct device *dev, struct ioctl_request *req) {
    /* Validate key pointer — in kernel this would be copy_from_user */
    if (!req->key) {
        fprintf(stderr, "ioctl: SET_CONFIG: NULL key pointer\n");
        return -1;
    }

    /* Check key length */
    size_t key_len = strlen(req->key);
    if (key_len == 0 || key_len >= MAX_KEY_LEN) {
        fprintf(stderr, "ioctl: SET_CONFIG: invalid key length %zu\n", key_len);
        return -1;
    }

    /* Look for existing entry to update */
    struct config_entry *entry = find_config(dev, req->key);
    if (!entry) {
        entry = find_free_slot(dev);
        if (!entry) {
            fprintf(stderr, "ioctl: SET_CONFIG: config table full\n");
            return -1;
        }
        strncpy(entry->key, req->key, MAX_KEY_LEN - 1);
        entry->key[MAX_KEY_LEN - 1] = '\0';
        entry->in_use = 1;
        dev->config_count++;
    }

    /* Copy value from request into config entry */
    size_t copy_len = req->value_len;
    if (copy_len >= MAX_VALUE_LEN) {
        copy_len = MAX_VALUE_LEN - 1;
    }
    memcpy(entry->value, req->value, copy_len);
    entry->value[copy_len] = '\0';

    return 0;
}

/*
 * Handle IOCTL_GET_CONFIG: retrieve a configuration value by key.
 */
static int handle_get_config(struct device *dev, struct ioctl_request *req) {
    if (!req->key) {
        fprintf(stderr, "ioctl: GET_CONFIG: NULL key pointer\n");
        return -1;
    }

    struct config_entry *entry = find_config(dev, req->key);
    if (!entry) {
        fprintf(stderr, "ioctl: GET_CONFIG: key '%s' not found\n", req->key);
        return -1;
    }

    /* Copy value to output buffer */
    size_t copy_len = strlen(entry->value);
    if (copy_len >= req->value_len) {
        copy_len = req->value_len - 1;
    }
    memcpy(req->value, entry->value, copy_len);
    req->value[copy_len] = '\0';

    return 0;
}

/*
 * Main ioctl dispatch function.
 * Validates the device pointer and request pointer, then dispatches
 * to the appropriate handler based on the command code.
 */
int handle_ioctl(struct device *dev, struct ioctl_request *req) {
    /* Validate device pointer */
    if (!dev) {
        fprintf(stderr, "ioctl: NULL device pointer\n");
        return -1;
    }

    /* Validate request pointer */
    if (!req) {
        fprintf(stderr, "ioctl: NULL request pointer\n");
        return -1;
    }

    /* Check device is initialized */
    if (!dev->initialized) {
        fprintf(stderr, "ioctl: device not initialized\n");
        return -1;
    }

    switch (req->cmd) {
    case IOCTL_SET_CONFIG:
        return handle_set_config(dev, req);

    case IOCTL_GET_CONFIG:
        return handle_get_config(dev, req);

    case IOCTL_RESET:
        memset(dev->config, 0, sizeof(dev->config));
        dev->config_count = 0;
        fprintf(stderr, "ioctl: device config reset\n");
        return 0;

    case IOCTL_GET_STATUS:
        printf("Status: %s, configs: %d, flags: 0x%lx\n",
               dev->initialized ? "UP" : "DOWN",
               dev->config_count, dev->flags);
        return 0;

    default:
        fprintf(stderr, "ioctl: unknown command %d\n", req->cmd);
        return -1;
    }
}
