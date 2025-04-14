package com.example.farmscurity

import android.app.NotificationChannel
import android.app.NotificationManager
import android.os.Build
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.annotation.RequiresApi
import androidx.compose.animation.core.tween
import androidx.compose.animation.slideInHorizontally
import androidx.compose.animation.slideOutHorizontally
import androidx.compose.runtime.Composable
import androidx.navigation.NavType
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import androidx.navigation.navArgument
import com.example.farmscurity.screen.Capture
import com.example.farmscurity.screen.Home
import com.example.farmscurity.ui.theme.FarmScurityTheme


class MainActivity : ComponentActivity() {
    @RequiresApi(Build.VERSION_CODES.N)
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        createNotificationChannel()
        setContent {
            FarmScurityTheme {
                Navigation()
            }
        }
    }

    private fun createNotificationChannel() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channel = NotificationChannel(
                "mqtt_channel_id",
                "Gerakan Terdeteksi!",
                NotificationManager.IMPORTANCE_HIGH
            ).apply {
                description = "Channel untuk notifikasi MQTT"
            }
            val notificationManager = getSystemService(NotificationManager::class.java)
            notificationManager.createNotificationChannel(channel)
        }
    }
}

@Composable
fun Navigation() {
    val navController = rememberNavController()

    NavHost(
        navController = navController,
        startDestination = "home",
        enterTransition = {
            slideInHorizontally(initialOffsetX = { it }, animationSpec = tween(300))
        },
        exitTransition = {
            slideOutHorizontally(targetOffsetX = { -it }, animationSpec = tween(300))
        }
    ) {
        composable("home") { Home(navController) }
        composable(
            route = "capture/{pictureId}/{isCapture}",
            arguments = listOf(
                navArgument("pictureId") { type = NavType.StringType },
                navArgument("isCapture") { type = NavType.StringType }
            )
        ) { backStackEntry ->
            val pictureId = backStackEntry.arguments?.getString("pictureId")
            val isCapture = backStackEntry.arguments?.getString("isCapture")?.toBoolean() ?: false

            Capture(
                navController = navController,
                pictureId = pictureId,
                isCapture = isCapture)
        }
    }
}
