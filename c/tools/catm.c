#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/stat.h>
#include <sys/mman.h>

// ----------------------------------------------------------------
static void emit(char* sol, char* eol, FILE* output_stream) {
	 size_t ntowrite = eol - sol;
     size_t nwritten = fwrite(sol, 1, ntowrite, output_stream);
	 if (nwritten != ntowrite) {
		perror("fwrite");
		exit(1);
	 }
}

// ----------------------------------------------------------------
static int do_stream(char* file_name) {
	FILE* output_stream = stdout;
	int fd = open(file_name, O_RDONLY);
	if (fd < 0) {
		perror("open");
		exit(1);
	}
	struct stat stat;
	if (fstat(fd, &stat) < 0) {
		perror("fstat");
		exit(1);
	}
	char* sof = mmap(NULL, (size_t)stat.st_size, PROT_READ, MAP_FILE|MAP_PRIVATE, fd, (off_t)0);
	if (sof == MAP_FAILED) {
		perror("mmap");
		exit(1);
	}
	char* eof = sof + stat.st_size;
	char* sol = sof;
	char* eol;
	char* p = sof;

	while (p < eof) {
		if (*p == '\n') {
			p++;
			eol = p;
			emit(sol, eol, output_stream);
			sol = p;
		} else {
			p++;
		}
	}

//	while (1) {
//		char* line = NULL;
//		size_t linecap = 0;
//		ssize_t linelen = getdelim(&line, &linecap, '\n', input_stream);
//		if (linelen <= 0) {
//			break;
//		}
//		fputs(line, output_stream);
//		free(line);
//	}

	if (close(fd) < 0) {
		perror("close");
		exit(1);
	}

	return 1;
}

// ================================================================
int main(int argc, char** argv) {
	int ok = 1;
	for (int argi = 1; argi < argc; argi++) {
	    ok = do_stream(argv[argi]);
	}
	return ok ? 0 : 1;
}
