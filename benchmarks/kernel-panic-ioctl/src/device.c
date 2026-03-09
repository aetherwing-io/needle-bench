#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "device.h"
#include "ioctl_handler.h"

struct device *device_create(const char *name) {
    struct device *dev = calloc(1, sizeof(struct device));
    if (!dev)
        return NULL;

    strncpy(dev->name, name, DEVICE_NAME_LEN - 1);
    dev->name[DEVICE_NAME_LEN - 1] = '\0';
    dev->config_count = 0;
    dev->initialized = 1;
    dev->flags = 0;

    /* Set some default configuration */
    struct ioctl_request req;
    req.cmd = IOCTL_SET_CONFIG;
    req.key = "mode";
    req.value = "normal";
    req.value_len = 6;
    handle_ioctl(dev, &req);

    req.key = "buffer_size";
    req.value = "4096";
    req.value_len = 4;
    handle_ioctl(dev, &req);

    return dev;
}

void device_destroy(struct device *dev) {
    if (dev) {
        dev->initialized = 0;
        free(dev);
    }
}

void device_print_info(struct device *dev) {
    printf("Device: %s\n", dev->name);
    printf("Initialized: %s\n", dev->initialized ? "yes" : "no");
    printf("Config entries: %d\n", dev->config_count);
    printf("Flags: 0x%lx\n", dev->flags);

    for (int i = 0; i < MAX_CONFIG_ENTRIES; i++) {
        if (dev->config[i].in_use) {
            printf("  %s = %s\n", dev->config[i].key, dev->config[i].value);
        }
    }
}

/* Run basic self-test with known-good inputs */
int device_self_test(struct device *dev) {
    int failures = 0;

    printf("Running device self-test...\n");

    /* Test 1: set and get a config value */
    struct ioctl_request set_req;
    set_req.cmd = IOCTL_SET_CONFIG;
    set_req.key = "test_key";
    set_req.value = "test_value";
    set_req.value_len = 10;

    if (handle_ioctl(dev, &set_req) != 0) {
        printf("  FAIL: could not set test_key\n");
        failures++;
    }

    struct ioctl_request get_req;
    get_req.cmd = IOCTL_GET_CONFIG;
    get_req.key = "test_key";
    char buf[256];
    get_req.value = buf;
    get_req.value_len = sizeof(buf);

    if (handle_ioctl(dev, &get_req) != 0) {
        printf("  FAIL: could not get test_key\n");
        failures++;
    } else if (strcmp(buf, "test_value") != 0) {
        printf("  FAIL: expected 'test_value', got '%s'\n", buf);
        failures++;
    }

    /* Test 2: update existing key */
    set_req.value = "updated";
    set_req.value_len = 7;
    if (handle_ioctl(dev, &set_req) != 0) {
        printf("  FAIL: could not update test_key\n");
        failures++;
    }

    if (handle_ioctl(dev, &get_req) != 0) {
        printf("  FAIL: could not get updated test_key\n");
        failures++;
    } else if (strcmp(buf, "updated") != 0) {
        printf("  FAIL: expected 'updated', got '%s'\n", buf);
        failures++;
    }

    /* Test 3: get nonexistent key */
    get_req.key = "nonexistent";
    if (handle_ioctl(dev, &get_req) == 0) {
        printf("  FAIL: getting nonexistent key should fail\n");
        failures++;
    }

    if (failures == 0) {
        printf("Self-test PASSED\n");
    } else {
        printf("Self-test FAILED: %d failures\n", failures);
    }

    return failures;
}

/* Fuzz test with crafted/adversarial inputs */
int device_fuzz_test(struct device *dev) {
    int failures = 0;

    printf("Running fuzz test with crafted inputs...\n");

    /* Test 1: NULL request pointer */
    printf("  Test: NULL request pointer... ");
    int ret = handle_ioctl(dev, NULL);
    if (ret == 0) {
        printf("FAIL (should have returned error)\n");
        failures++;
    } else {
        printf("OK (returned %d)\n", ret);
    }

    /* Test 2: NULL key in set request */
    printf("  Test: NULL key in set request... ");
    struct ioctl_request bad_req;
    bad_req.cmd = IOCTL_SET_CONFIG;
    bad_req.key = NULL;
    bad_req.value = "test";
    bad_req.value_len = 4;
    ret = handle_ioctl(dev, &bad_req);
    if (ret == 0) {
        printf("FAIL (should have returned error)\n");
        failures++;
    } else {
        printf("OK (returned %d)\n", ret);
    }

    /* Test 3: NULL value in set request */
    printf("  Test: NULL value in set request... ");
    bad_req.key = "valid_key";
    bad_req.value = NULL;
    bad_req.value_len = 100;
    ret = handle_ioctl(dev, &bad_req);
    if (ret == 0) {
        printf("FAIL (should have returned error)\n");
        failures++;
    } else {
        printf("OK (returned %d)\n", ret);
    }

    /* Test 4: NULL value buffer in get request */
    printf("  Test: NULL value buffer in get request... ");
    bad_req.cmd = IOCTL_GET_CONFIG;
    bad_req.key = "mode";
    bad_req.value = NULL;
    bad_req.value_len = 0;
    ret = handle_ioctl(dev, &bad_req);
    if (ret == 0) {
        printf("FAIL (should have returned error)\n");
        failures++;
    } else {
        printf("OK (returned %d)\n", ret);
    }

    /* Test 5: Invalid ioctl command */
    printf("  Test: invalid ioctl command... ");
    bad_req.cmd = 9999;
    bad_req.key = "test";
    bad_req.value = "test";
    bad_req.value_len = 4;
    ret = handle_ioctl(dev, &bad_req);
    if (ret == 0) {
        printf("FAIL (should have returned error)\n");
        failures++;
    } else {
        printf("OK (returned %d)\n", ret);
    }

    /* Test 6: NULL device pointer */
    printf("  Test: NULL device pointer... ");
    struct ioctl_request valid_req;
    valid_req.cmd = IOCTL_GET_CONFIG;
    valid_req.key = "mode";
    char buf[256];
    valid_req.value = buf;
    valid_req.value_len = sizeof(buf);
    ret = handle_ioctl(NULL, &valid_req);
    if (ret == 0) {
        printf("FAIL (should have returned error)\n");
        failures++;
    } else {
        printf("OK (returned %d)\n", ret);
    }

    /* Test 7: Oversized value_len */
    printf("  Test: oversized value_len... ");
    bad_req.cmd = IOCTL_SET_CONFIG;
    bad_req.key = "oversize_test";
    bad_req.value = "small";
    bad_req.value_len = 999999;  /* Way larger than actual string */
    ret = handle_ioctl(dev, &bad_req);
    if (ret == 0) {
        printf("FAIL (should have returned error or truncated safely)\n");
        failures++;
    } else {
        printf("OK (returned %d)\n", ret);
    }

    if (failures == 0) {
        printf("Fuzz test PASSED\n");
    } else {
        printf("Fuzz test FAILED: %d failures\n", failures);
    }

    return failures;
}
