#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "ioctl_handler.h"
#include "device.h"

static void print_usage(void) {
    fprintf(stderr, "Usage: devctl <command> [args...]\n");
    fprintf(stderr, "\n");
    fprintf(stderr, "Commands:\n");
    fprintf(stderr, "  info                    Show device info\n");
    fprintf(stderr, "  set-config <key> <val>  Set device configuration\n");
    fprintf(stderr, "  get-config <key>        Get device configuration\n");
    fprintf(stderr, "  test                    Run self-test\n");
    fprintf(stderr, "  fuzz                    Run fuzz test with crafted inputs\n");
}

int main(int argc, char *argv[]) {
    if (argc < 2) {
        print_usage();
        return 1;
    }

    struct device *dev = device_create("vdev0");
    if (!dev) {
        fprintf(stderr, "Error: failed to create device\n");
        return 1;
    }

    int ret = 0;

    if (strcmp(argv[1], "info") == 0) {
        device_print_info(dev);
    } else if (strcmp(argv[1], "set-config") == 0) {
        if (argc < 4) {
            fprintf(stderr, "Usage: devctl set-config <key> <value>\n");
            ret = 1;
        } else {
            struct ioctl_request req;
            req.cmd = IOCTL_SET_CONFIG;
            req.key = argv[2];
            req.value = argv[3];
            req.value_len = strlen(argv[3]);
            ret = handle_ioctl(dev, &req);
        }
    } else if (strcmp(argv[1], "get-config") == 0) {
        if (argc < 3) {
            fprintf(stderr, "Usage: devctl get-config <key>\n");
            ret = 1;
        } else {
            struct ioctl_request req;
            req.cmd = IOCTL_GET_CONFIG;
            req.key = argv[2];
            char buf[256];
            req.value = buf;
            req.value_len = sizeof(buf);
            ret = handle_ioctl(dev, &req);
            if (ret == 0) {
                printf("%s = %s\n", argv[2], buf);
            }
        }
    } else if (strcmp(argv[1], "test") == 0) {
        ret = device_self_test(dev);
    } else if (strcmp(argv[1], "fuzz") == 0) {
        ret = device_fuzz_test(dev);
    } else {
        fprintf(stderr, "Unknown command: %s\n", argv[1]);
        print_usage();
        ret = 1;
    }

    device_destroy(dev);
    return ret;
}
