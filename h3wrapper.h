#ifndef H3WRAPPER_H_
#define H3WRAPPER_H_

#include <h3lib/h3lib.h>

int h3lib_create_bucket(H3_Handle handle, char *bucketName);
int h3lib_purge_bucket(H3_Handle handle, char *bucketName);
int h3lib_write_object(H3_Handle handle, char *bucketName, char *objectName, void *data, uint64_t size, uint64_t offset);
int h3lib_read_dummy_object(H3_Handle handle, char *bucketName, char *objectName);
int h3lib_delete_object(H3_Handle handle, char *bucketName, char *objectName);

int h3lib_write_object_from_file(H3_Handle handle, char *bucketName, char *objectName, char *filename);
int h3lib_read_object_to_file(H3_Handle handle, char *bucketName, char *objectName, char *filename);

#endif
