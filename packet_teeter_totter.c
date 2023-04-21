#include <arpa/inet.h>
#include <netinet/ip.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <unistd.h>

#define BUF_SIZE 4096
#define PORT_A 12345
#define PORT_B 12346

int is_whitelisted(uint32_t ip_addr) {
    // Implement your whitelist checking logic here.
    // For this example, let's say we have a single whitelisted IP address: 192.168.1.100
    uint32_t whitelisted_ip = inet_addr("192.168.1.100");
    return (ip_addr == whitelisted_ip);
}

int main() {
    int raw_socket = socket(AF_INET, SOCK_RAW, IPPROTO_TCP);
    if (raw_socket < 0) {
        perror("Failed to create raw socket");
        exit(EXIT_FAILURE);
    }

    struct sockaddr_in dest_addr;
    memset(&dest_addr, 0, sizeof(dest_addr));
    dest_addr.sin_family = AF_INET;
    dest_addr.sin_addr.s_addr = inet_addr("127.0.0.1");

    char buffer[BUF_SIZE];
    while (1) {
        memset(buffer, 0, BUF_SIZE);

        struct sockaddr_in src_addr;
        socklen_t addr_len = sizeof(src_addr);
        int recv_len = recvfrom(raw_socket, buffer, BUF_SIZE, 0, (struct sockaddr *)&src_addr, &addr_len);
        if (recv_len < 0) {
            perror("Failed to receive packet");
            continue;
        }

        struct iphdr *ip_header = (struct iphdr *)buffer;

        if (is_whitelisted(ip_header->saddr)) {
            dest_addr.sin_port = htons(PORT_A);
        } else {
            dest_addr.sin_port = htons(PORT_B);
        }

        int send_len = sendto(raw_socket, buffer, recv_len, 0, (struct sockaddr *)&dest_addr, sizeof(dest_addr));
        if (send_len < 0) {
            perror("Failed to send packet");
            continue;
        }
    }

    close(raw_socket);
    return 0;
}
