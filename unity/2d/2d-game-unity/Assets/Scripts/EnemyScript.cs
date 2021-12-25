using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using Example;
public class EnemyScript : MonoBehaviour
{
    private bool hasSpawn;
    private MoveScript moveScript;
    private WeaponScript[] weapons;
    private Collider2D coliderComponent;
    private SpriteRenderer rendererComponent;
    void Awake()
    {
        // 初始化時 獲取組件應用
        weapons = GetComponentsInChildren<WeaponScript>();
        moveScript = GetComponent<MoveScript>();
        coliderComponent = GetComponent<Collider2D>();
        rendererComponent = GetComponent<SpriteRenderer>();
    }
    void Start()
    {
        hasSpawn = false;

        // 禁用所有功能組件
        coliderComponent.enabled = false;
        moveScript.enabled = false;
        foreach (WeaponScript weapon in weapons)
        {
            weapon.enabled = false;
        }
    }


    void Update()
    {
        if (hasSpawn == false)
        {
            if (rendererComponent.IsVisibleFrom(Camera.main))//進入到屏幕
            {
                Spawn();
            }
        }
        else
        {
            foreach (var weapon in weapons)
            {
                // 自動開火
                if (weapon != null && weapon.CanAttack)
                {
                    weapon.Attack(true);
                }
            }
            // 離開屏幕 銷毀對象
            if (rendererComponent.IsVisibleFrom(Camera.main) == false)
            {
                Destroy(gameObject);
            }
        }
    }
    private void Spawn()
    {
        hasSpawn = true;

        // 啓用各種組件
        coliderComponent.enabled = true;
        moveScript.enabled = true;
        foreach (WeaponScript weapon in weapons)
        {
            weapon.enabled = true;
        }
    }
}
