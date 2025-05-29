package com.example.farmscurity.screen

import android.Manifest
import android.content.Intent
import android.os.Build
import android.widget.Toast
import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.Image
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.AlertDialog
import androidx.compose.material.Divider
import androidx.compose.material.OutlinedButton
import androidx.compose.material.Switch
import androidx.compose.material.SwitchDefaults
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Delete
import androidx.compose.material.icons.filled.Info
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.ModalBottomSheet
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.material3.rememberModalBottomSheetState
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.SideEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.shadow
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.lifecycle.viewmodel.compose.viewModel
import androidx.navigation.NavController
import com.example.farmscurity.service.MQTTServiceNotification
import com.example.farmscurity.R
import com.example.farmscurity.model.ApiServiceViewModel
import com.example.farmscurity.model.History
import com.example.farmscurity.model.Sensor
import com.example.farmscurity.ui.theme.MyTypography
import com.google.accompanist.permissions.ExperimentalPermissionsApi
import com.google.accompanist.permissions.isGranted
import com.google.accompanist.permissions.rememberPermissionState
import com.google.accompanist.systemuicontroller.rememberSystemUiController
import kotlinx.coroutines.launch

@Composable
fun SensorScreen(viewModel: ApiServiceViewModel = viewModel()) {
    val sensorItems by viewModel.sensorList.collectAsState()

    Text(text = "Sensor", fontSize = 20.sp, fontWeight = FontWeight.Bold)
    Spacer(modifier = Modifier.height(8.dp))
    if (sensorItems.isEmpty()) {
        Text(
            text = "Tidak ada data",
            modifier = Modifier.fillMaxWidth(),
            textAlign = TextAlign.Center
        )
    } else {
        Card(
            modifier = Modifier
                .fillMaxWidth()
                .padding(vertical = 4.dp)
                .shadow(
                    elevation = 2.dp,
                    shape = RoundedCornerShape(8.dp),
                    ambientColor = Color.Black.copy(alpha = 0.75f),
                    spotColor = Color.Black.copy(alpha = 0.75f)
                ),
            shape = RoundedCornerShape(8.dp),
            colors = CardDefaults.cardColors(containerColor = Color.White)
        )
        {
            Column {
                Spacer(modifier = Modifier.height(12.dp))
                Column {
                    sensorItems.forEachIndexed { index, sensor ->
                        Column(
                            modifier = Modifier.padding(horizontal = 16.dp)
                        ) {

                            SensorItem(sensor)
                            Spacer(modifier = Modifier.height(12.dp))
                            if (index < sensorItems.lastIndex) {
                                Divider(
                                    color = Color(0xFFDFDFDF),
                                    thickness = 1.dp,
                                )
                                Spacer(modifier = Modifier.height(12.dp))
                            }
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun HistoryScreen(viewModel: ApiServiceViewModel = viewModel(), navController: NavController) {
    val historyItems by viewModel.historyList.collectAsState()

    Text(text = "Riwayat", fontSize = 20.sp, fontWeight = FontWeight.Bold)
    Spacer(modifier = Modifier.height(8.dp))
    if (historyItems.isEmpty()) {
        Text(
            text = "Tidak ada data",
            modifier = Modifier.fillMaxWidth(),
            textAlign = TextAlign.Center
        )
    } else {
        Card(
            modifier = Modifier
                .fillMaxWidth()
                .padding(vertical = 4.dp)
                .shadow(
                    elevation = 2.dp,
                    shape = RoundedCornerShape(8.dp),
                    ambientColor = Color.Black.copy(alpha = 0.75f),
                    spotColor = Color.Black.copy(alpha = 0.75f)
                ),
            shape = RoundedCornerShape(8.dp),
            colors = CardDefaults.cardColors(containerColor = Color.White)
        )
        {
            Column {
                Button(
                    onClick = { viewModel.deleteHistory() },
                    shape = RoundedCornerShape(8.dp),
                    colors = ButtonDefaults.buttonColors(containerColor = Color(0xFFFF2B2B)),
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(16.dp)
                ) {
                    Icon(imageVector = Icons.Default.Delete, contentDescription = "Hapus", tint = Color.White)
                    Spacer(modifier = Modifier.width(8.dp))
                    Text("Hapus Semua Riwayat", color = Color.White, style = MyTypography.bodySmall,)
                }
                Spacer(modifier = Modifier.height(12.dp))
                Column {
                    historyItems.forEachIndexed { index, history ->
                        if (!history.isRead) {
                            AlertDialog(
                                onDismissRequest = {  },
                                title = {
                                    Box(
                                        modifier = Modifier
                                            .fillMaxWidth(),
                                        contentAlignment = Alignment.Center
                                    ){
                                        Icon(
                                            imageVector = Icons.Default.Info,
                                            contentDescription = null,
                                            tint = Color.Red,
                                            modifier = Modifier
                                                .size(70.dp)

                                        )
                                    }
                                },
                                text = {
                                    Column(
                                        modifier = Modifier
                                            .fillMaxWidth(),
                                        horizontalAlignment = Alignment.CenterHorizontally
                                    ) {
                                        Text(
                                            text = "Gerakan Terdeteksi‼️",
                                            fontWeight = FontWeight.Bold,
                                            fontSize = 18.sp,
                                            color = Color.Black,
                                            textAlign = TextAlign.Center
                                        )
                                        Spacer(modifier = Modifier.height(8.dp))
                                        Text(
                                            text = history.description,
                                            fontSize = 14.sp,
                                            color = Color.Gray,
                                            textAlign = TextAlign.Center
                                        )
                                    }
                                },
                                confirmButton = {
                                    Box(
                                        modifier = Modifier
                                            .fillMaxWidth()
                                            .background(
                                                color = Color(0xFF45556C),
                                                shape = RoundedCornerShape(10.dp)
                                            )
                                            .clickable {
                                                viewModel.setIsRead(history.id)
                                                navController.navigate("capture/${history.pictureId}/false")
                                            }
                                            .padding(vertical = 12.dp, horizontal = 24.dp)
                                    ) {
                                        Text(
                                            "Lihat Info",
                                            modifier = Modifier.align(Alignment.Center),
                                            style = MyTypography.titleMedium,
                                            color = Color.White
                                        )
                                    }
                                },
                                dismissButton = {}
                            )
                        }
                        Column(
                            modifier = Modifier.padding(horizontal = 16.dp)
                        ) {

                            HistoryItem(history, navController)
                            Spacer(modifier = Modifier.height(12.dp))
                            if (index < historyItems.lastIndex) {
                                Divider(
                                    color = Color(0xFFDFDFDF),
                                    thickness = 1.dp,
                                )
                                Spacer(modifier = Modifier.height(12.dp))
                            }
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun SensorItem(sensor: Sensor){
    Row (
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.SpaceBetween,
        modifier = Modifier.fillMaxWidth()
    ) {
        Row (
            verticalAlignment = Alignment.CenterVertically
        ){
            Image(
                painter = painterResource(id = R.drawable.pir_dark),
                contentDescription = "PIR Dark",
                modifier = Modifier.size(23.dp)
            )
            Spacer(modifier = Modifier.width(8.dp))
            Text(
                text = "Passive Infrared Sensor ${sensor.id}",
                style = MyTypography.bodySmall
            )
        }
        Box(
            modifier = Modifier
                .size(16.dp)
                .clip(CircleShape)
                .background(if (sensor.isActive) Color(0xFF4CAF50) else Color.Red)
        )
    }
}

@Composable
fun SensorItemChecklist(sensor: Sensor, viewModel: ApiServiceViewModel){
    val coroutineScope = rememberCoroutineScope()
    var isChecked by remember { mutableStateOf(sensor.isActive) }
    var previousState by remember { mutableStateOf(sensor.isActive) }

    LaunchedEffect(sensor.id, sensor.isActive) {
        isChecked = sensor.isActive
    }

    Row (
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.SpaceBetween,
        modifier = Modifier.fillMaxWidth()
    ) {
        Row (
            verticalAlignment = Alignment.CenterVertically
        ) {
            Image(
                painter = painterResource(id = R.drawable.pir_dark),
                contentDescription = "PIR Dark",
                modifier = Modifier.size(23.dp)
            )
            Spacer(modifier = Modifier.width(8.dp))
            Text(
                text = "Passive Infrared Sensor ${sensor.id}",
                style = MyTypography.bodySmall
            )
        }

        Switch(
            checked = isChecked,
            onCheckedChange = { newValue ->
                previousState = isChecked
                isChecked = newValue

                coroutineScope.launch {
                    viewModel.setIsActive(sensor.id, newValue) { success ->
                        if (!success) {
                            isChecked = previousState
                            viewModel.fetchSensor()
                        }
                    }
                }
            },
            colors = SwitchDefaults.colors(
                checkedThumbColor = Color.White,
                checkedTrackColor = Color(0xFF4CAF50),
                uncheckedThumbColor = Color.White,
                uncheckedTrackColor = Color(0xFFBDBDBD)
            )
        )

    }
}

@Composable
fun HistoryItem(history: History, navController: NavController) {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .clickable { navController.navigate("capture/${history.pictureId}/false") }
    ) {
        Text(text = history.operation, style = MyTypography.labelSmall)
        Spacer(modifier = Modifier.padding(top = 4.dp))
        Text(text = history.description, style = MyTypography.bodySmall)
        Spacer(modifier = Modifier.padding(top = 8.dp))
        Text(text = history.createdAt, color = Color(0xFF9D9D9D), style = MyTypography.bodySmall)
    }
}

@Composable
fun InfoBox(viewModel: ApiServiceViewModel = viewModel()) {
    val sensorItems by viewModel.sensorList.collectAsState()
    val historyItem by viewModel.historyList.collectAsState()

    Surface(
        shape = RoundedCornerShape(10.dp),
        modifier = Modifier.height(80.dp),
    ) {
        if (sensorItems.isEmpty() && historyItem.isEmpty()) {
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .background(Color(0xFF45556C)),
                contentAlignment = Alignment.Center
            ) {
                Text(
                    text = "Tidak ada data",
                    modifier = Modifier.fillMaxWidth(),
                    textAlign = TextAlign.Center,
                    color = Color.White
                )
            }
        } else {
            Box(
                modifier = Modifier
                    .fillMaxWidth()
                    .background(Color(0xFF45556C)),
                contentAlignment = Alignment.Center
            ) {
                Row(
                    verticalAlignment = Alignment.CenterVertically,
                    horizontalArrangement = Arrangement.SpaceEvenly,
                    modifier = Modifier.fillMaxWidth()
                ) {
                    Row(
                        verticalAlignment = Alignment.CenterVertically) {
                        Image(
                            painter = painterResource(id = R.drawable.pir_light),
                            contentDescription = "PIR Icon",
                            modifier = Modifier.size(30.dp)
                        )
                        Spacer(modifier = Modifier.width(10.dp))
                        Column {
                            Text(
                                text = sensorItems.count { it.isActive }.toString(),
                                fontSize = 28.sp,
                                fontWeight = FontWeight.Bold,
                                color = Color.White
                            )
                            Text(
                                text = "Sensor Aktif",
                                fontSize = 14.sp,
                                fontWeight = FontWeight.Normal,
                                color = Color.White
                            )
                        }
                    }
                    Box(
                        modifier = Modifier
                            .height(50.dp)
                            .width(1.dp)
                            .background(Color(0xFF617796))
                    )
                    Column(
                        horizontalAlignment = Alignment.CenterHorizontally,
                        verticalArrangement = Arrangement.Center
                    ) {
                        Text(
                            text = "Aktivitas Terakhir",
                            fontSize = 12.sp,
                            color = Color.White,
                            fontWeight = FontWeight.Bold
                        )
                        Text(
                            text = historyItem.firstOrNull()?.createdAt ?: "Tidak Ada Data",
                            fontSize = 12.sp,
                            color = Color.White,
                            fontWeight = FontWeight.Normal
                        )
                    }
                }
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun BottomSheet(viewModel: ApiServiceViewModel = viewModel()) {
    val sheetState = rememberModalBottomSheetState()
    var showBottomSheet by remember { mutableStateOf(false) }
    val sensorItems by viewModel.sensorList.collectAsState()

    Column(
        modifier = Modifier.fillMaxSize(),
        verticalArrangement = Arrangement.Center,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Box(
            modifier = Modifier
                .fillMaxWidth()
                .background(color = Color(0xFF45556C), shape = RoundedCornerShape(10.dp))
                .clickable { showBottomSheet = true }
                .padding(vertical = 12.dp, horizontal = 24.dp)
        ) {
            Text("Ubah Sensor",
                modifier = Modifier.align(Alignment.Center),
                style = MyTypography.titleMedium,
                color = Color.White)
        }

        if (showBottomSheet) {
            ModalBottomSheet(
                onDismissRequest = { showBottomSheet = false },
                sheetState = sheetState
            ) {
                Column (
                    modifier = Modifier
                        .background(color = Color.White)
                        .padding(horizontal = 24.dp)
                ) {
                    Text(text = "Ubah Sensor", style = MyTypography.bodyLarge)
                    Spacer(modifier = Modifier.height(24.dp))
                    sensorItems.forEachIndexed{index, sensor ->
                        SensorItemChecklist(sensor, viewModel)
                        Spacer(modifier = Modifier.height(12.dp))
                        if (index < sensorItems.lastIndex) {
                            Divider(
                                color = Color(0xFFDFDFDF),
                                thickness = 1.dp,
                            )
                            Spacer(modifier = Modifier.height(12.dp))
                        }
                    }
                }

            }
        }
    }
}

@OptIn(ExperimentalPermissionsApi::class)
@Composable
fun Home(navController: NavController, viewModel: ApiServiceViewModel = viewModel()) {
    val context = LocalContext.current
    var toastMessage by remember { mutableStateOf<String?>(null) }
    val postNotificationPermission = rememberPermissionState(
        permission = Manifest.permission.POST_NOTIFICATIONS
    )

    var isAlramActive by remember { mutableStateOf(false) }

    val systemUiController = rememberSystemUiController()
    val backgroundColor = Color(0xFFF5F7FF)

    val capture by viewModel.capture.collectAsState()

    SideEffect {
        systemUiController.setStatusBarColor(color = backgroundColor)
    }

    LaunchedEffect(Unit) {
        if (!postNotificationPermission.status.isGranted) {
            postNotificationPermission.launchPermissionRequest()
        }

        val intent = Intent(context, MQTTServiceNotification::class.java)
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            context.startForegroundService(intent)
        } else {
            context.startService(intent)
        }
    }

    toastMessage?.let { message ->
        LaunchedEffect(message) {
            Toast.makeText(context, message, Toast.LENGTH_SHORT).show()
            toastMessage = null
        }
    }

    MaterialTheme(typography = MyTypography) {
        Surface(
            modifier = Modifier
                .fillMaxSize()
                .padding(top = 26.dp),
            color = backgroundColor
        ) {
            LazyColumn(
                modifier = Modifier
                    .fillMaxSize()
                    .padding(start = 32.dp, end = 32.dp, bottom = 12.dp)
            ) {
                item { Spacer(modifier = Modifier.height(32.dp)) }
                item {
                    Row(
                        verticalAlignment = Alignment.CenterVertically,
                        horizontalArrangement = Arrangement.SpaceBetween,
                        modifier = Modifier.fillMaxWidth()
                    ) {
                        Text(
                            text = "Halo, User",
                            style = MyTypography.bodyLarge
                        )

                    }
                }

                item { Spacer(modifier = Modifier.height(32.dp)) }

                // Info Box
                item {
                    InfoBox()
                }
                item { Spacer(modifier = Modifier.height(8.dp)) }
                item {
                    val coroutineScope = rememberCoroutineScope()
                    val context = LocalContext.current

                    Box(
                        modifier = Modifier
                            .fillMaxWidth()
                            .background(
                                color = Color(0xFF45556C),
                                shape = RoundedCornerShape(10.dp)
                            )
                            .clickable {
                                toastMessage = "Menangkap gambar"

                                coroutineScope.launch {
                                    val result = viewModel.capture()

                                    if (result) {
                                        val pictureId = viewModel.capture.value.firstOrNull()?.data
                                        if (pictureId != null) {
                                            navController.navigate("capture/$pictureId/true")
                                            toastMessage = "Berhasil mengambil gambar"
                                        } else {
                                            toastMessage = "ID gambar tidak ditemukan"
                                        }
                                    } else {
                                        toastMessage = "Gagal mengambil gambar"
                                    }
                                }
                            }
                            .padding(vertical = 12.dp, horizontal = 24.dp)
                    ) {
                        Text(
                            "Tangkap Gambar",
                            modifier = Modifier.align(Alignment.Center),
                            style = MyTypography.titleMedium,
                            color = Color.White
                        )
                    }

                }

                item { Spacer(modifier = Modifier.height(4.dp)) }
                item {
                    Button(
                        onClick = {
                            isAlramActive = !isAlramActive
                            viewModel.alarm(isActive = isAlramActive) {onResult ->
                                if (!onResult){
                                    isAlramActive = !isAlramActive
                                    toastMessage = "Gagal menyalakan Alarm"
                                }
                            }
                        },
                        colors = ButtonDefaults.buttonColors(
                            containerColor = backgroundColor
                        ),
                        shape = RoundedCornerShape(10.dp),
                        border = BorderStroke(2.dp, Color(0xFF45556C)),
                        modifier = Modifier
                            .fillMaxWidth()

                    ) {
                        Text(
                            text = if (!isAlramActive) "Nyalakan Alarm" else "Matikan Alarm",
                            style = MyTypography.titleMedium,
                            color = Color(0xFF45556C)
                        )
                    }


                }
                item { Spacer(modifier = Modifier.height(25.dp)) }
                item { SensorScreen() }
                item { Spacer(modifier = Modifier.height(8.dp)) }
                item { BottomSheet() }
                item { Spacer(modifier = Modifier.height(25.dp)) }
                item { HistoryScreen(navController = navController) }
            }
        }
    }
}