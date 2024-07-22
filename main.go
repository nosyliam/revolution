package main

// #include <stdio.h>
// #include <semaphore.h>
// int destroy_semaphore()
// {
//     const char *sem_name = "/RobloxPlayerUniq";
//     if (sem_unlink(sem_name) == -1) // Attempt to destroy the semaphore
//     {
//         return 1;
//     }
//     return 0;
// }
import "C"

func main() {
	C.destroy_semaphore()
}
