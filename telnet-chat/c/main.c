#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <pthread.h>

#define BUF_SIZE 100
#define NAME_TIP "Please enter your name: "
#define COMMAND_TIP "> "
#define BASE_IP "127.0.0.1"
#define PORT 1234

enum ROLE
{
      GUEST, MEMBER, ADMIN
};

struct User;

struct User
{
   char username[30];
   char display[30];
   char password[18];
   enum ROLE role;
};

struct Client
{
   int sock;
   struct User user;
};

struct SockClient
{
   int serverSock;
   struct sockaddr_in clnt_addr;
   socklen_t clnt_addr_size;
};

void createClient(void *ptr);
void createSock();

struct User users[10];
struct SockClient sock;

int main()
{
   createSock();
   while (1)
   {
      int client_sock = accept(sock.serverSock, (struct sockaddr *)&sock.clnt_addr,
         &sock.clnt_addr_size);
      pthread_t client_thread;
      if (client_sock > 0) {
         int ret_thread = pthread_create(&client_thread,NULL,(void *)&createClient,(void *)&client_sock);
      }
   }

   //关闭套接字
   close(sock.serverSock);
   return 0;
}

// 创建 sock
void createSock()
{
   //创建套接字
   sock.serverSock = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
   //将套接字和IP、端口绑定
   struct sockaddr_in serv_addr;
   memset(&serv_addr, 0, sizeof(serv_addr));
   serv_addr.sin_family = AF_INET; //使用IPv4地址
   serv_addr.sin_addr.s_addr = inet_addr(BASE_IP);
   serv_addr.sin_port = htons(PORT);
   bind(sock.serverSock, (struct sockaddr *)&serv_addr, sizeof(serv_addr));
   //进入监听状态，等待用户发起请求
   listen(sock.serverSock, 20);
   struct sockaddr_in clnt_addr;
   sock.clnt_addr_size = sizeof(clnt_addr);
}

// 接收一个连接
void createClient(void *ptr)
{
   struct Client client;
   client.sock = *(int*)ptr;

   char buffer[BUF_SIZE] = ""; //缓冲区
   char rcvMsg[200] = "";
   while (1)
   {
      // 设置用户名
      if (strcmp(client.user.username, "") == 0)
      {
         write(client.sock, NAME_TIP, sizeof(NAME_TIP));   // 发送数据
         int strLen = read(client.sock, buffer, BUF_SIZE); //接收客户端发来的数据
         if (strLen <= 0) { // 客户端断开连接
            close(client.sock);
            break;
         }
         strcpy(client.user.username, buffer);
         strcpy(rcvMsg, "您的用户名为: ");
         strcat(rcvMsg, client.user.username);
         write(client.sock, rcvMsg, sizeof(rcvMsg)); // 发送数据
      }
      else
      {
         write(client.sock, COMMAND_TIP, sizeof(COMMAND_TIP)); // 发送数据
         int strLen = read(client.sock, buffer, BUF_SIZE); //接收客户端发来的数据
         if (strLen <= 0) {
            close(client.sock);
            break;
         }
         strcpy(rcvMsg, "已收到消息: ");
         strcat(rcvMsg, buffer);
         write(client.sock, rcvMsg, sizeof(rcvMsg)); // 发送数据
      }
      memset(buffer, 0, BUF_SIZE); //重置缓冲区
      memset(rcvMsg, 0, 200); //重置缓冲区
   }
}