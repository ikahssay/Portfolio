package com.application.cs260final;

import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.Toast;

public class Login extends AppCompatActivity {
    Button login;


    @Override
    protected void onCreate(Bundle savedInstanceState) {
          super.onCreate(savedInstanceState);
          setContentView(R.layout.activity_login_screen);
        login = (Button) findViewById(R.id.btnLOGIN);
        login.setOnClickListener(new View.OnClickListener() {

            //When button is clicked, run code inside
            public void onClick(View v) {
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
               // Toast.makeText(Login.this, "Log New Activity", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(Login.this, HomeScreenActivity.class);

                //Start Activity on a new screen that's encoded in the Choose_Time_or_Distance_Activity class file.
                startActivity(intent);
            }
        });

      }
}
