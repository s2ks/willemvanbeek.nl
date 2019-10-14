#include <stdio.h>
#include <stdlib.h>
#include <string.h>

extern char **environ;

#define BUF_SIZE 1024

int main(void)
{
	char *buf, *input = NULL, *end;

	buf = malloc(BUF_SIZE);

	for(int i = 1; fgets(buf, BUF_SIZE, stdin) != NULL; i++) {
		input = end = realloc(input, i * BUF_SIZE);
		end += strnlen(input, BUF_SIZE);
		snprintf(end, BUF_SIZE, "%s", buf);
	}


	for(int i = 0; environ[i] != NULL; i++) {
		printf("%s\n", environ[i]);
	}

	puts(buf);

	free(buf);
	free(input);
	return 0;
}
