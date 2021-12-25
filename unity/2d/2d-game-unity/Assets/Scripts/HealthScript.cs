using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class HealthScript : MonoBehaviour
{
    //  定義總血量
    public int hp = 1;
    // 是否是敵人
    public bool isEnemy = true;
    // 此函數用於計算傷害
    public void Damage(int damageCount)
    {
        hp -= damageCount;
        if (hp <= 0)
        {
            // 粒子效果
            SpecialEffectsHelper.Instance.Explosion(transform.position);
            // 死亡
            Destroy(gameObject);
        }
    }
    // 和觸發器發生了碰撞
    void OnTriggerEnter2D(Collider2D otherCollider)
    {
        // 獲取發射物腳本
        var shot = otherCollider.gameObject.GetComponent<ShotScript>();
        if (shot != null) // 只有和發射物腳本碰撞才執行
        {
            if (shot.isEnemyShot != isEnemy)    // 只有被對立面的發射物擊中才執行
            {
                // 計算傷害
                Damage(shot.damage);

                // 銷毀發射物
                Destroy(shot.gameObject); // 記住應該始終以 GameObject 爲目標，否則你將只刪除腳本
            }
        }
    }
}
