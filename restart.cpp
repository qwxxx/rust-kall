#include <stdlib.h>
#include <signal.h>
#include <chrono>
#include <thread>

int main(){
    char buf[512];
    FILE *cmd_pipe = popen("ps -A | grep sspar | awk '{print $1}'", "r");

    fgets(buf, 512, cmd_pipe);
    pid_t pid = strtoul(buf, NULL, 10);
    pclose( cmd_pipe );  
    system("./sspar &");
    // std::this_thread::sleep_for(std::chrono::milliseconds(1000));
    kill(pid, SIGKILL);
    return 0;
}