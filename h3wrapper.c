#include <sys/stat.h>
#include <fcntl.h>

#include "h3wrapper.h"

#define TRUE 1;
#define FALSE 0;

int h3lib_create_bucket(H3_Handle handle, char *bucketName) {
    H3_Auth token = {0};

    return H3_CreateBucket(handle, &token, bucketName);
}

int h3lib_purge_bucket(H3_Handle handle, char *bucketName) {
    H3_Auth token = {0};

    return H3_PurgeBucket(handle, &token, bucketName);
}

int h3lib_write_object(H3_Handle handle, char *bucketName, char *objectName, void *data, uint64_t size, uint64_t offset) {
    H3_Auth token = {0};

    return H3_WriteObject(handle, &token, bucketName, objectName, data, size, offset);
}

int h3lib_read_dummy_object(H3_Handle handle, char *bucketName, char *objectName) {
    H3_Auth token = {0};
    size_t size;

    return H3_ReadDummyObject(handle, &token, bucketName, objectName, &size);
}

int h3lib_delete_object(H3_Handle handle, char *bucketName, char *objectName) {
    H3_Auth token = {0};

    return H3_DeleteObject(handle, &token, bucketName, objectName);
}

int h3lib_write_object_from_file(H3_Handle handle, char *bucketName, char *objectName, char *filename) {
    off_t offset = 0;
    uint32_t userId = 0;

    H3_Auth auth;
    auth.userId = userId;

    struct stat st;
    if (stat(filename, &st) == -1) {
        return FALSE;
    }
    size_t size = st.st_size;
    int fd = open(filename, O_RDONLY);
    if (fd == -1) {
        return FALSE;
    }

    H3_Status return_value = H3_WriteObjectFromFile(handle, &auth, bucketName, objectName, fd, size, offset);
    if (return_value != H3_SUCCESS && return_value != H3_CONTINUE) {
        close(fd);
        return FALSE;
    }
    close(fd);

    return TRUE;
}

int h3lib_read_object_to_file(H3_Handle handle, char *bucketName, char *objectName, char *filename) {
    off_t offset = 0;
    size_t size = 0;
    uint32_t userId = 0;

    H3_Auth auth;
    auth.userId = userId;

    int fd = open(filename, O_CREAT|O_WRONLY|O_TRUNC, 0644);
    if (fd == -1) {
        return FALSE;
    }

    H3_Status return_value = H3_ReadObjectToFile(handle, &auth, bucketName, objectName, offset, fd, &size);
    if (return_value != H3_SUCCESS && return_value != H3_CONTINUE) {
        close(fd);
        return FALSE;
    }
    close(fd);

    return TRUE;
}
