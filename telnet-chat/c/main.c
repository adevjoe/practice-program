#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <pthread.h>

#define BUF_SIZE 200
#define HELP_TIP "\
Chat Room.\n\n\
Usage:\n\
  /command [args...]\n\n\
Commands:\n\
   register <username> <password>    create an account\n\
   login <username> <password>       login with your username and password\n\
   logout                            logout your account\n\
   join <room_id>                    join a room\n\
   leave <room_id>                   leave a room\n\
   quit                              quit current room\n\
"
#define COMMAND_TIP "> "
#define UNKNOWN_COMMAND_TIP "Unknown Command.\n"
#define BASE_IP "127.0.0.1"
#define PORT 1234

enum ROLE
{
      GUEST, MEMBER, ADMIN
};

struct User;

struct User
{
   int ID;
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

struct Room
{
   int ID;
   char name[30];
   int users[100];
};

struct SockClient
{
   int serverSock;
   struct sockaddr_in clnt_addr;
   socklen_t clnt_addr_size;
};

void createClient(void *ptr);
void createSock();

// user interface
struct User createUser(char username[30], char password[18]);
struct User getUserByUsername(char username[30]);
struct User getUserByID(int id);
int checkUser(char username[30], char password[18]);
int setUserRole(int id, enum ROLE role);

// room interface
struct Room createRoom(char name[30]);
struct Room listRoom();
int joinRoom(int userID, int roomID);

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

void handlerMessage(struct Client client, char *msg) {
   if (msg[0] == '/')
   {
      write(client.sock, msg, sizeof(msg));
   } else
   {
      write(client.sock, UNKNOWN_COMMAND_TIP, sizeof(UNKNOWN_COMMAND_TIP));
   }
}

// 接收一个连接
void createClient(void *ptr)
{
   struct Client client;
   client.sock = *(int*)ptr;

   char buffer[BUF_SIZE] = ""; //缓冲区
   char rcvMsg[BUF_SIZE] = "";
   write(client.sock, HELP_TIP, sizeof(HELP_TIP)); // 发送数据
   while (1)
   {
      write(client.sock, COMMAND_TIP, sizeof(COMMAND_TIP)); // 发送数据
      int strLen = read(client.sock, buffer, BUF_SIZE); //接收客户端发来的数据
      if (strLen <= 0) {
         close(client.sock);
         break;
      }
      // handler message
      strcat(rcvMsg, buffer);
      handlerMessage(client, rcvMsg);
      memset(buffer, 0, BUF_SIZE);
      memset(rcvMsg, 0, BUF_SIZE);
   }
}