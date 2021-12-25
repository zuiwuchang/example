using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class PlayerScript : MonoBehaviour
{
    // 定義一個表示速度的屬性，設置爲 public 可以在 unity 可視化窗口中調整值
    public Vector2 speed = new Vector2(50, 50);
    // Start is called before the first frame update

    // 2 - Store the movement and the component
    private Vector2 movement;// 速度
    private Rigidbody2D rigidbodyComponent; // 緩存 剛體對象
    void Start()
    {

    }

    // Update is called once per frame
    void Update()
    {
        // 獲取用戶輸入
        float inputX = Input.GetAxis("Horizontal");
        float inputY = Input.GetAxis("Vertical");

        // 依據輸入設置 當前速度
        movement = new Vector2(
          speed.x * inputX,
          speed.y * inputY);

        bool shoot = Input.GetButtonDown("Fire1");
        if (!shoot)
        {
            shoot = Input.GetButtonDown("Fire2");
        }

        if (shoot)
        {
            WeaponScript weapon = GetComponent<WeaponScript>();
            if (weapon != null)
            {
                // 因爲是玩家在使用所以傳入 false 參數
                weapon.Attack(false);
            }
        }

        // 確保 Player 不在攝像機之外
        var dist = (transform.position - Camera.main.transform.position).z;
        var leftBorder = Camera.main.ViewportToWorldPoint(new Vector3(0, 0, dist)).x;
        var rightBorder = Camera.main.ViewportToWorldPoint(new Vector3(1, 0, dist)).x;
        var topBorder = Camera.main.ViewportToWorldPoint(new Vector3(0, 0, dist)).y;
        var bottomBorder = Camera.main.ViewportToWorldPoint(new Vector3(0, 1, dist)).y;

        transform.position = new Vector3(
          Mathf.Clamp(transform.position.x, leftBorder, rightBorder),
          Mathf.Clamp(transform.position.y, topBorder, bottomBorder),
          transform.position.z
        );
    }
    void FixedUpdate()
    {
        // 獲取組件並
        if (rigidbodyComponent == null)
        {
            rigidbodyComponent = GetComponent<Rigidbody2D>();
        }

        // 設置速度移動剛體
        rigidbodyComponent.velocity = movement;
    }

    void OnCollisionEnter2D(Collision2D collision)
    {
        bool damagePlayer = false;

        // 獲取碰撞體上 掛載的 EnemyScript
        EnemyScript enemy = collision.gameObject.GetComponent<EnemyScript>();
        if (enemy != null)
        {
            // 殺手敵人
            HealthScript enemyHealth = enemy.GetComponent<HealthScript>();
            if (enemyHealth != null)
            {
                enemyHealth.Damage(enemyHealth.hp);
            }

            damagePlayer = true;
        }

        // 敵人對玩家造成傷害
        if (damagePlayer)
        {
            HealthScript playerHealth = this.GetComponent<HealthScript>();
            if (playerHealth != null)
            {
                playerHealth.Damage(1);
            }
        }
    }

    void OnDestroy()
    {
        // 銷毀 遊戲結束
        var gameOver = FindObjectOfType<GameOverScript>();
        gameOver.ShowButtons();
    }
}
