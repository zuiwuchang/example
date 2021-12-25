using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class SpecialEffectsHelper : MonoBehaviour
{
    // 單例以供其它腳本使用
    public static SpecialEffectsHelper Instance;
    // 煙霧 prefab
    public ParticleSystem smokeEffect;
    // 火焰 prefab
    public ParticleSystem fireEffect;
    // Start is called before the first frame update
    void Awake()
    {
        // 此腳本應該只被添加到一個單獨實例上
        if (Instance != null)
        {
            Debug.LogError("Multiple instances of SpecialEffectsHelper!");
        }
        // 初始化單例
        Instance = this;
    }

    // 創建爆炸效果
    public void Explosion(Vector3 position)
    {
        // 實例化煙霧
        instantiate(smokeEffect, position);

        // 實例化火焰
        instantiate(fireEffect, position);
    }
    // 實例化 粒子系統
    private ParticleSystem instantiate(ParticleSystem prefab, Vector3 position)
    {
        ParticleSystem newParticleSystem = Instantiate(
          prefab,
          position,
          Quaternion.identity
        ) as ParticleSystem;

        // Make sure it will be destroyed
        Destroy(
          newParticleSystem.gameObject,
          newParticleSystem.main.startLifetimeMultiplier
        );

        return newParticleSystem;
    }
}
