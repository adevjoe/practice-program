#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netinet/in.h>

#define BUF_SIZE 100
#define NAME_TIP "Please enter your name: "
#define COMMAND_TIP "> "

struct User;

struct User
{
   char name[100];
};

struct Client
{
   int sock;
   struct User user;
};

// 接收一个连接
void createClient(int sock)
{
   struct Client client;
   client.sock = sock;

   char buffer[BUF_SIZE] = ""; //缓冲区
   char rcvMsg[200] = "";
   for (;;)
   {
      // 设置用户名
      if (strcmp(client.user.name, "") == 0)
      {
         write(client.sock, NAME_TIP, sizeof(NAME_TIP));   // 发送数据
         int strLen = read(client.sock, buffer, BUF_SIZE); //接收客户端发来的数据
         strcpy(client.user.name, buffer);
         strcpy(rcvMsg, "您的用户名为: ");
         strcat(rcvMsg, client.user.name);
         write(client.sock, rcvMsg, sizeof(rcvMsg)); // 发送数据
      }
      else
      {
         write(client.sock, COMMAND_TIP, sizeof(COMMAND_TIP)); // 发送数据
         int strLen = read(client.sock, buffer, BUF_SIZE);     //接收客户端发来的数据
         strcpy(rcvMsg, "已收到消息: ");
         strcat(rcvMsg, buffer);
         write(client.sock, rcvMsg, sizeof(rcvMsg)); // 发送数据
      }
      memset(buffer, 0, BUF_SIZE); //重置缓冲区
      memset(rcvMsg, 0, 200);      //重置缓冲区
   }
}

int main()
{

   //创建套接字
   int serv_sock = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
   //将套接字和IP、端口绑定
   struct sockaddr_in serv_addr;
   memset(&serv_addr, 0, sizeof(serv_addr));           //每个字节都用0填充
   serv_addr.sin_family = AF_INET;                     //使用IPv4地址
   serv_addr.sin_addr.s_addr = inet_addr("127.0.0.1"); //具体的IP地址
   serv_addr.sin_port = htons(1234);                   //端口
   bind(serv_sock, (struct sockaddr *)&serv_addr, sizeof(serv_addr));
   //进入监听状态，等待用户发起请求
   listen(serv_sock, 20);
   //接收客户端请求
   struct sockaddr_in clnt_addr;
   socklen_t clnt_addr_size = sizeof(clnt_addr);
   for (;;)
   {
      int client_sock = accept(serv_sock, (struct sockaddr *)&clnt_addr, &clnt_addr_size);
      createClient(client_sock);
   }

   //关闭套接字
   // close(serv_sock);
   return 0;
}