// Copyright (c) 2020 FORTH-ICS
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
