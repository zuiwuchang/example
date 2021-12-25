using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class MoveScript : MonoBehaviour
{
    // 移動速度
    public Vector2 speed = new Vector2(10, 10);

    // 移動方向
    public Vector2 direction = new Vector2(-1, 0);

    private Vector2 movement;
    private Rigidbody2D rigidbodyComponent;

    void Update()
    {
        // 再每幀都重設速度以便讓敵人保持一直移動
        movement = new Vector2(
                speed.x * direction.x,
                speed.y * direction.y);
    }
    void FixedUpdate()
    {
        if (rigidbodyComponent == null)
        {
            rigidbodyComponent = GetComponent<Rigidbody2D>();
        }

        // 將速度應用到剛體
        rigidbodyComponent.velocity = movement;
    }
}
