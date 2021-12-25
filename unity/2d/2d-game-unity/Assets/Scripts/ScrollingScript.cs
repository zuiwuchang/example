using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using System.Linq; // 爲 List 提供了 OrderBy 函數
using Example; // 引入 IsVisibleFrom 擴展
public class ScrollingScript : MonoBehaviour
{
    // 滾動速度
    public Vector2 speed = new Vector2(2, 2);
    //  滾動方向
    public Vector2 direction = new Vector2(-1, 0);
    // 是否連接到攝像機
    public bool isLinkedToCamera = false;
    // 是否是無限滾動的
    public bool isLooping = false;
    // 對於無限滾動背景 存儲其渲染子節點
    private List<SpriteRenderer> backgroundPart;

    void Start()
    {
        // 對於無限滾動的執行一些初始化
        if (isLooping)
        {
            // 獲取所有存在渲染器的子節點
            backgroundPart = new List<SpriteRenderer>();
            for (int i = 0; i < transform.childCount; i++)
            {
                Transform child = transform.GetChild(i);
                SpriteRenderer r = child.GetComponent<SpriteRenderer>();
                // 值添加可見節點
                if (r != null)
                {
                    backgroundPart.Add(r);
                }
            }

            // 按照子節點位置 從左到右排序
            backgroundPart = backgroundPart.OrderBy(
              t => t.transform.position.x
            ).ToList();
        }
    }

    void Update()
    {
        // 移動
        Vector3 movement = new Vector3(
          speed.x * direction.x,
          speed.y * direction.y,
          0);
        movement *= Time.deltaTime;
        transform.Translate(movement);

        // 對於連接了攝像機的節點 同時移動攝像機
        if (isLinkedToCamera)
        {
            Camera.main.transform.Translate(movement);
        }
        if (isLooping)
        {
            loopUpdate();
        }
    }
    private void loopUpdate()
    {
        // 獲取列表最左邊的元素
        SpriteRenderer firstChild = backgroundPart.FirstOrDefault();
        if (firstChild == null)
        {
            return;
        }
        // 檢查節點是否已經在攝像機之前 因爲 IsVisibleFrom 函數比較耗資源
        if (firstChild.transform.position.x < Camera.main.transform.position.x)
        {
            if (!firstChild.IsVisibleFrom(Camera.main))// 沒在鏡頭內
            {
                // 獲取最後一個節點的位置
                SpriteRenderer lastChild = backgroundPart.LastOrDefault();

                Vector3 lastPosition = lastChild.transform.position;
                Vector3 lastSize = (lastChild.bounds.max - lastChild.bounds.min);

                // 將第一個節點移到到最後一個節點之後，目前只處理了水平滾動
                firstChild.transform.position = new Vector3(lastPosition.x + lastSize.x,
                    firstChild.transform.position.y, firstChild.transform.position.z);

                // 移動子節點儀表
                backgroundPart.Remove(firstChild);
                backgroundPart.Add(firstChild);
            }
        }
    }
}
