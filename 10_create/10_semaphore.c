#include <semaphore.h>

sem_t *go_sem_open(const char *name) {
	return sem_open(name, O_CREAT, 0644, 1);
}
