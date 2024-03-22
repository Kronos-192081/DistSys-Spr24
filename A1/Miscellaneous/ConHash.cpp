#include <iostream>
#include <vector>
#include <ctime>
#include <map>
#include <string>

using namespace std;

const int mod = 1e5+7;

struct Node{
    bool occ;
    string name;
};

class ConHash{
    vector<Node> hash_d;
    int size;
    int vs;
    int len = 0;
    map<string, int> all_server;
    map<string, int> server_id;

    int get_serv_hash(int i, int j)
    {
        long long val = 1LL*i*i + 1LL*j*j + 2*j + 25;
        val = val%size;
        return (int)val;
    }

    int get_cli_hash(int i)
    {
        long long val = 1LL*i*i + 1LL*2*i + 17;
        val = val%size;
        return (int)val;
    }
    
    public:

    ConHash(int m, int k): size{m}, vs{k} {
        hash_d.resize(m);
        for(int i = 0; i<m; i++){
            hash_d[i] = {false, ""};
        }
    }

    ~ConHash(){}

    int add(vector<int> ids, vector<string> names){

        if(ids.size() != names.size()) return 0;

        if((len+ids.size())*vs >= size) return 0;
        len+=ids.size();
        
        for(int i = 0; i<(int)ids.size(); i++)
        {
            all_server[names[i]] = 1;
            server_id[names[i]] = ids[i];
            for(int j = 0; j<vs; j++)
            {
                int hash = get_serv_hash(ids[i], j);
                while(hash_d[hash].occ){
                    hash = (hash + 1)%size;
                }
                hash_d[hash] = {true, names[i]};  
            }
        }
        return 1;
    }

    void get_config()
    {
        for(int i = 0; i<size; i++)
        {
            cout << "Index: " << i << " ";
            cout << "Status: " << hash_d[i].occ << " ";
            cout << "Server: " << hash_d[i].name << endl;
        }
    }

    int add(int id, string name){
        if((len+1)*vs >= size) return 0;
        len++;
        all_server[name] = 1;
        server_id[name] = id;
        for(int j = 0; j<vs; j++)
        {
            int hash = get_serv_hash(id, j);
            while(hash_d[hash].occ)
            {
                hash = (hash + 1)%size;
            }
            hash_d[hash] = {true, name};
        }

        return 1;
    }

    int rem(string name){
        if(all_server.find(name) == all_server.end()) return 0;

        for(int j = 0; j<vs; j++)
        {
            int hash = get_serv_hash(server_id[name], j);
            while(hash_d[hash].name != name) { hash = (hash + 1)%size; }
            hash_d[hash] = {false, ""};
        }
        all_server.erase(name);
        server_id.erase(name);
        len--;
        return 1;
    }

    string get_serv(int id)
    {
        if(len == 0) return "No Server Allocable";
        int hash = get_cli_hash(id);
        // cout << id << " : " << hash << endl;
        hash = (hash + 1)%size;
        while(!hash_d[hash].occ)
        {
            hash = (hash + 1)%size;
        }

        return hash_d[hash].name;
    }
};

int main()
{
    srand(time(0));
    ConHash c(10, 3);
   
    vector<int> v = {rand()%mod, rand()%mod};
    vector<string> s = {"Server 1", "Server 2"};

    c.add(v, s);
    // c.get_config();

    int cnt = 5;
    
    cout << "\n####### Testing Round 1 #######" << endl;
    while(cnt--)
    {
        int id = rand()%mod;
        cout << id << " : ";
        cout << c.get_serv(id) << endl;
    }

    c.add(rand()%mod, "Server 3");
    cout << "\nAdded Server 3\n";
    cnt = 5;

    cout << "\n####### Testing Round 2 #######" << endl;


    while(cnt--)
    {
        int id = rand()%mod;
        cout << id << " : ";
        cout << c.get_serv(id) << endl;
    }

    c.rem("Server 1");
    cout << "\nRemoved server 1\n";

    cnt = 5;
    cout << "\n####### Testing Round 3 #######" << endl;

    while(cnt--)
    {
        int id = rand()%mod;
        cout << id << " : ";
        cout << c.get_serv(id) << endl;
    }
    
    return 0;
}
