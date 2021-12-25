using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class ShotScript : MonoBehaviour
{
    // 定義將造成多少傷害
    public int damage = 1;
    // 標記是敵人還是玩家發射的彈藥
    public bool isEnemyShot = false;
    void Start()
    {
        // 延遲20秒後銷毀射擊物，即彈藥被發射後的最大存活時間爲20秒
        Destroy(gameObject, 20);
    }
}
