using UnityEngine;
using UnityEngine.UI;
using UnityEngine.SceneManagement;
public class GameOverScript : MonoBehaviour
{
    private Button[] buttons;

    void Awake()
    {
        // 獲取所有按鈕
        buttons = GetComponentsInChildren<Button>();

        // 隱藏它們直到玩家死亡
        HideButtons();
    }

    public void HideButtons()
    {
        foreach (var b in buttons)
        {
            b.gameObject.SetActive(false);
        }
    }

    public void ShowButtons()
    {
        foreach (var b in buttons)
        {
            b.gameObject.SetActive(true);
        }
    }

    public void ExitToMenu()
    {
        // 重新加載菜單場景
        SceneManager.LoadScene("Menu");
    }

    public void RestartGame()
    {
        // 重新加載遊戲
        SceneManager.LoadScene("Main");
    }
}
