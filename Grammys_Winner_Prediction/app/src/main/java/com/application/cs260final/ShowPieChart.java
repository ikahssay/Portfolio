package com.application.cs260final;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.ImageView;
import android.widget.Toast;

import androidx.appcompat.app.AppCompatActivity;

public class ShowPieChart extends AppCompatActivity {
    ImageView home_button;
    ImageView search_button;
    ImageView notification_button;
    ImageView profile_button;

    @Override
    protected void onCreate(Bundle savedInstanceState){
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_voted_chart);


        /*
        //TODO: CANT FIND NEXT PAGE TO THIS ACTION!!!
        //1. CREATE AN ACTION LINKED TO THE VOTING RADIO BUTTON
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        show_pie_chart_button = (RadioButton) findViewById(R.id.show_pie_chart_button);
        show_pie_chart_button.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(ShowPieChart.this, "View Voting Breakdown", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(ShowPieChart.this, Login.class); //TODO: FIX!!

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });
*/

        // ************************************* MENU BUTTON FUNCTIONALITIES BELOW **************************************** //

        //2. CREATE AN ACTION LINKED TO THE PREDICTION BUTTON (AT THE MENU BAR BELOW)
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        //TODO: ALL MENU BUTTONS ARE GROUPED AS ONE. CANNOT GET SEARCH BUTTON'S ID
        home_button = (ImageView) findViewById(R.id.imageHome);
        home_button.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(ShowPieChart.this, "Going to Home Page", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(ShowPieChart.this, HomeScreenActivity.class);

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });

        //4. CREATE AN ACTION LINKED TO THE SEARCH BUTTON (AT THE MENU BAR BELOW)
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        //TODO: ALL MENU BUTTONS ARE GROUPED AS ONE. CANNOT GET SEARCH BUTTON'S ID
        search_button = (ImageView) findViewById(R.id.imageSearch);
        search_button.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(ShowPieChart.this, "Going to Search Page", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(ShowPieChart.this, Search.class);

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });

        //6. CREATE AN ACTION LINKED TO THE NOTIFICATIONS BUTTON (AT THE MENU BAR BELOW)
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        //TODO: ALL MENU BUTTONS ARE GROUPED AS ONE. CANNOT GET SEARCH BUTTON'S ID
        notification_button = (ImageView) findViewById(R.id.imageBell);
        notification_button.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(ShowPieChart.this, "Going to Notifications Page", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(ShowPieChart.this, Notifications.class);

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });

        //7. CREATE AN ACTION LINKED TO THE PROFILE BUTTON (AT THE MENU BAR BELOW)
        //-> setOnClickListener is listening if button is ever clicked. If it is, all the actions in this function run.
        //TODO: ALL MENU BUTTONS ARE GROUPED AS ONE. CANNOT GET SEARCH BUTTON'S ID
        profile_button = (ImageView) findViewById(R.id.imageProfileButton);
        profile_button.setOnClickListener(new View.OnClickListener(){

            //When button is clicked, run code inside
            public void onClick(View v) {
                //TODO: CAN REMOVE THIS AFTER EVERYTHING WORKS CORRECTLY
                //specifies what class its in, what the text says, how long the text will be there, and if we should show it or not.
                Toast.makeText(ShowPieChart.this, "Going to Profile Page", Toast.LENGTH_SHORT).show();

                //Links the action of pressing button = going to next layout.
                Intent intent = new Intent(ShowPieChart.this, Profile.class);

                //Start Activity on a new screen that's encoded in the activity_home_screen_click_on_post class file.
                startActivity(intent);
            }
        });
    }
}