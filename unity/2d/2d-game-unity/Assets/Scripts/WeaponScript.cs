using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class WeaponScript : MonoBehaviour
{
    // 定義要實例化的發射物
    public Transform shotPrefab;
    // 兩次射擊的 冷卻時間
    public float shootingRate = 0.25f;
    // 還有多久冷卻
    private float shootCooldown;
    // Start is called before the first frame update
    void Start()
    {
        // 設置冷卻
        shootCooldown = 0f;
    }

    // Update is called once per frame
    void Update()
    {
        if (shootCooldown > 0)// 如果武器過熱
        {
            // 拖進時間，以使武器冷卻
            shootCooldown -= Time.deltaTime;
        }
    }
    // 此函數用於武器進行射擊
    public void Attack(bool isEnemy)
    {
        if (CanAttack)//判斷是否冷卻
        {
            shootCooldown = shootingRate;// 射擊後設置冷卻時間

            // 創建一個發射物實例
            var shotTransform = Instantiate(shotPrefab) as Transform;

            // 修改彈藥方向和武器方向一致
            shotTransform.position = transform.position;

            // 設置武器是否是敵人在使用
            ShotScript shot = shotTransform.gameObject.GetComponent<ShotScript>();
            if (shot != null)
            {
                shot.isEnemyShot = isEnemy;
            }

            // 使用彈藥方向始終和移動方向一致
            MoveScript move = shotTransform.gameObject.GetComponent<MoveScript>();
            if (move != null)
            {
                move.direction = this.transform.right; // 在 2D 空間中，精靈的朝向是精靈的左側
            }
        }
    }
    // 返回是否冷卻
    public bool CanAttack
    {
        get
        {
            return shootCooldown <= 0f;
        }
    }
}
